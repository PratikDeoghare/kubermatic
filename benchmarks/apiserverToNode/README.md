#### Introduction 
There are two scripts to run. 

**NOTE:** Run `./run-qperf-client-in-pod.sh` first it will tell you when and with what args to run `./run-qperf-server-on-node.sh`. 

#### Assumptions

- Nodes are running ubuntu
- SSH agent is enabled
- apiserver is runnning in alpine based pod
- Expose strategy is tunneling

#### Script ./run-qperf-client-in-pod.sh

This script takes

- seed-cluster kubeconfig
- user-cluster kubeconfig
- seed-cluster namespace corresponding to the user-cluster

Like this

`
./run-qperf-client-in-pod.sh <seed-kubeconfig> <user-kubeconfig> <seed-namespace>
`

1. The script then dumps info about current cluster state.

2. Picks a node to run qperf server and prompts you to run the script to run qperf server on the chosen
   node (`./run-qperf-server-on-node.sh`).
3. Once the server is running on the node, you allow this script to proceed.
4. Then the script picks an apiserver pod and installs qperf in it.
5. Applies appropriate iptable rules.
6. Runs qperf client to calculate tcp_bw and tcp_lat.

#### Script ./run-qperf-server-on-node.sh

It takes ip address of the node to run qperf server on.

Run this script when prompted by the other script and with argument suggested by that script.
