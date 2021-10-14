#!/usr/bin/env bash

set -ex

seedconfig=$1
userconfig=$2
ns=$3
bid=$RANDOM$RANDOM
function seed() {
  kubectl --kubeconfig=$seedconfig --namespace=$ns $@
}

function user() {
  kubectl --kubeconfig=$userconfig $@
}

function log() {
  echo
  echo "[BENCH-$bid] " $@
  echo
}

log "user-cluster state"
user get all -A -owide
user get nodes -A -owide

log "seed-cluster state"
seed get all -owide

nodeStatus=$(user get nodes -ojsonpath='{.items[0].status.addresses}')
nodeExternalIP=$(echo $nodeStatus | jq '.[] | select(.type | contains("ExternalIP")) | .address' -r)
nodeInternalIP=$(echo $nodeStatus | jq '.[] | select(.type | contains("InternalIP")) | .address' -r)
log "InternalIP:" $nodeInternalIP
log "ExternalIP:" $nodeExternalIP

right=$(echo $nodeInternalIP | cut -d. -f3-4)
destIP=10.254.$right
log "DestinationIP:" $destIP

echo "Now go ahead and run this: "
read -p "bash ./run-qperf-server-on-node.sh $nodeExternalIP"

# setup
apiserverPod=$(seed get pods -l "app=apiserver" -oname | head -n1)

log "APIServer Pod:" $apiserverPod

log "installing qperf"
seed exec $apiserverPod -c openvpn-client -it -- /bin/bash <<EOF
  apk add --no-cache --repository http://dl-3.alpinelinux.org/alpine/edge/testing/ qperf==0.4.11-r0
EOF

log "modifying iptables:"
seed exec $apiserverPod -c openvpn-client -it -- /bin/bash <<EOF
  set -ex

  iptables -t nat -A node-access-dnat -d $nodeInternalIP/32 -p tcp -m tcp --dport 19766 -j DNAT --to-destination $destIP:19766
  iptables -t nat -A node-access-dnat -d $nodeInternalIP/32 -p tcp -m tcp --dport 19765 -j DNAT --to-destination $destIP:19765
  iptables -t nat -A node-access-dnat -d $nodeExternalIP/32 -p tcp -m tcp --dport 19766 -j DNAT --to-destination $destIP:19766
  iptables -t nat -A node-access-dnat -d $nodeExternalIP/32 -p tcp -m tcp --dport 19765 -j DNAT --to-destination $destIP:19765

  iptables-save
EOF

log "running benchmark..."
seed exec $apiserverPod -c openvpn-client -it -- /bin/bash <<EOF
  qperf -vvs -ip 19766 $nodeExternalIP tcp_bw tcp_lat
EOF
