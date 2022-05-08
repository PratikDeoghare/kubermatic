//go:build dualstack

/*
Copyright 2022 The Kubermatic Kubernetes Platform contributors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dualstack

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/kubermatic/machine-controller/pkg/cloudprovider/util"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	netutils "k8s.io/utils/net"
)

const (
	kubeSystem = "kube-system"
)

var (
	userconfig          string
	ipFamily            string
	skipNodes           bool
	skipHostNetworkPods bool
)

func init() {
	flag.StringVar(&userconfig, "userconfig", "", "path to kubeconfig of usercluster")
	flag.StringVar(&ipFamily, "ipFamily", "IPv4", "IP family")
	flag.BoolVar(&skipNodes, "skipNodes", true, "Set false to test nodes")
	flag.BoolVar(&skipHostNetworkPods, "skipHostNetworkPods", true, "Set false to test pods in host network")
}

func TestClusterIPFamily(t *testing.T) {
	// based on https://kubernetes.io/docs/tasks/network/validate-dual-stack/
	if userconfig == "" {
		t.Logf("kubeconfig for usercluster not provided, test passes vacuously.")
		t.Logf("to run against ready usercluster use following command")
		t.Logf("go test ./pkg/test/dualstack/dualstack -v -race -tags dualstack -timeout 30m -run TestClusterIPFamily -args --userconfig <USERCLUSTER KUBECONFIG> --ipFamily <IP FAMILY>")
		return
	}

	config, err := clientcmd.BuildConfigFromFlags("", userconfig)
	if err != nil {
		t.Fatalf("failed to build config: %s", err)
	}

	userclusterClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("failed to create usercluster client: %s", err)
	}

	testUserCluster(t, userclusterClient, util.IPFamily(ipFamily), skipNodes, skipHostNetworkPods)
}

func testUserCluster(t *testing.T, userclusterClient *kubernetes.Clientset, ipFamily util.IPFamily, skipNodes, skipHostNetworkPods bool) {
	t.Logf("testing with IP family: %q", ipFamily)
	ctx := context.Background()

	// validate nodes
	if skipNodes {
		t.Log("skipping validation for nodes")
	} else {
		nodes, err := userclusterClient.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
		if err != nil {
			t.Fatal(err)
		}

		for _, node := range nodes.Items {
			var addrs []string
			for _, addr := range node.Status.Addresses {
				fmt.Println(addr)
				if addr.Type == v1.NodeHostName {
					continue
				}
				addrs = append(addrs, addr.Address)
			}
			validate(t, node.Name, ipFamily, addrs)
		}

		for _, node := range nodes.Items {
			validate(t, node.Name, ipFamily, node.Spec.PodCIDRs)
		}
	}

	// validate pods
	{
		pods, err := userclusterClient.CoreV1().Pods(kubeSystem).List(ctx,
			metav1.ListOptions{})
		if err != nil {
			t.Fatal(err)
		}

		for _, pod := range pods.Items {
			if pod.Spec.HostNetwork && skipHostNetworkPods {
				t.Logf("skipping host network pod")
				continue
			}
			var podAddrs []string
			for _, addr := range pod.Status.PodIPs {
				podAddrs = append(podAddrs, addr.IP)
			}
			validate(t, pod.Name, ipFamily, podAddrs)
		}
	}

	// validate svc
	{
		svcs, err := userclusterClient.CoreV1().Services(kubeSystem).List(ctx, metav1.ListOptions{})
		if err != nil {
			t.Fatal(err)
		}

		for _, svc := range svcs.Items {
			switch *svc.Spec.IPFamilyPolicy {
			case v1.IPFamilyPolicySingleStack:
				if ipFamily == util.DualStack {
					t.Logf("skipping %q test for %q because IP family policy is %q", ipFamily, svc.Name, *svc.Spec.IPFamilyPolicy)
					continue
				}
			case v1.IPFamilyPolicyPreferDualStack, v1.IPFamilyPolicyRequireDualStack:
			}

			switch svc.Spec.Type {
			case v1.ServiceTypeClusterIP:
				validate(t, svc.Name, ipFamily, svc.Spec.ClusterIPs)
			case v1.ServiceTypeNodePort:
			case v1.ServiceTypeExternalName:
			case v1.ServiceTypeLoadBalancer:
				validate(t, svc.Name, ipFamily, svc.Spec.ClusterIPs)
				validate(t, svc.Name, ipFamily, svc.Spec.ExternalIPs)
			}
		}
	}

	// validate egress connectivity
	switch ipFamily {
	case util.IPv4, util.Unspecified:
		validateEgressConnectivity(t, userclusterClient, 4)
	case util.IPv6:
		validateEgressConnectivity(t, userclusterClient, 6)
	case util.DualStack:
		validateEgressConnectivity(t, userclusterClient, 4)
		validateEgressConnectivity(t, userclusterClient, 6)
	}
}

func validateEgressConnectivity(t *testing.T, userclusterClient *kubernetes.Clientset, ipVersion int) {
	t.Log("validating", fmt.Sprintf("egress-validator-%d", ipVersion))
	ns := "default"

	pod, err := userclusterClient.CoreV1().Pods(ns).Create(context.Background(), egressValidatorPod(ipVersion), metav1.CreateOptions{})
	if err != nil {
		t.Errorf("failed to create pod: %v", err)
		return
	}

	defer func() {
		err := userclusterClient.CoreV1().Pods(ns).Delete(context.Background(), pod.Name, metav1.DeleteOptions{})
		if err != nil {
			t.Errorf("failed to cleanup: %v", err)
		}
	}()

	err = wait.Poll(10*time.Second, 2*time.Minute, func() (bool, error) {
		p, err := userclusterClient.CoreV1().Pods(ns).Get(context.Background(), pod.Name, metav1.GetOptions{})
		if err != nil {
			t.Logf("failed to get pod: %v", err)
			return false, nil
		}

		if err := checkPodHealth(p); err != nil {
			t.Logf("pod not healthy: %v", err)
			return false, nil
		}

		return true, nil
	})

	if err != nil {
		t.Errorf("pod never became healthy: %v", err)
	}
}

func checkPodHealth(pod *v1.Pod) error {
	if pod.Status.Phase != v1.PodRunning {
		return fmt.Errorf("pod %q has phase not running: %s", pod.Name, pod.Status.Phase)
	}

	for _, c := range pod.Status.Conditions {
		if c.Type == v1.PodReady {
			if c.Status != v1.ConditionTrue {
				return fmt.Errorf("pod %q not ready", pod.Name)
			}
		} else if c.Type == v1.ContainersReady {
			if c.Status != v1.ConditionTrue {
				return fmt.Errorf("container not ready for pod %q", pod.Name)
			}
		}
	}

	return nil
}

func validate(t *testing.T, name string, ipFamily util.IPFamily, addrs []string) {
	fmt.Println("validating", name, addrs)
	if !all(ipFamily, addrs) {
		t.Errorf("not all addresses in %s are in IP family %q for %s", addrs, ipFamily, name)
	}
}

func all(ipFamily util.IPFamily, addrs []string) bool {
	// We convert all the IPs to CIDR notation so that we can simply use CIDR
	// validation functions everywhere instead of checking which function
	// to use every time.
	// Actual length of the mask doesn't matter, so it is set to 0.
	for i, addr := range addrs {
		if !strings.Contains(addr, "/") {
			addrs[i] = fmt.Sprintf("%s/0", addr)
		}
	}

	switch ipFamily {
	case util.IPv4, util.Unspecified:
		for _, addr := range addrs {
			if !netutils.IsIPv4CIDRString(addr) {
				return false
			}
		}
	case util.IPv6:
		for _, addr := range addrs {
			if !netutils.IsIPv6CIDRString(addr) {
				return false
			}
		}
	case util.DualStack:
		ok, err := netutils.IsDualStackCIDRStrings(addrs)
		return err == nil && ok
	default:
		return false
	}

	return true
}

func egressValidatorPod(ipVersion int) *v1.Pod {
	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: fmt.Sprintf("egress-validator-%d", ipVersion),
		},
		Spec: v1.PodSpec{
			Volumes:        nil,
			InitContainers: nil,
			Containers: []v1.Container{
				{
					Name:  fmt.Sprintf("egress-validator-%d-container", ipVersion),
					Image: "docker.io/byrnedo/alpine-curl:0.1.8",
					Command: []string{
						"/bin/ash",
						"-c",
						"sleep 1000000000",
					},
					LivenessProbe: &v1.Probe{
						ProbeHandler: v1.ProbeHandler{
							Exec: &v1.ExecAction{
								Command: []string{
									"curl",
									"-sS",
									fmt.Sprintf("--ipv%d", ipVersion),
									"--fail",
									"--connect-timeout",
									"5",
									"-o",
									"/dev/null",
									fmt.Sprintf("v%d.ident.me", ipVersion),
								},
							},
						},
						TimeoutSeconds: 7,
					},
					ReadinessProbe: &v1.Probe{
						ProbeHandler: v1.ProbeHandler{
							Exec: &v1.ExecAction{
								Command: []string{
									"curl",
									"-sS",
									fmt.Sprintf("--ipv%d", ipVersion),
									"--fail",
									"--connect-timeout",
									"5",
									"-o",
									"/dev/null",
									fmt.Sprintf("v%d.ident.me", ipVersion),
								},
							},
						},
						TimeoutSeconds: 7,
					},
				},
			},
			HostNetwork: false,
		},
	}
}
