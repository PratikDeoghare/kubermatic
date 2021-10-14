#!/usr/bin/env bash

qperfServerIP=$1 # ip of node where qperf server is running

ssh ubuntu@$qperfServerIP 'bash -s' <<EOF
  set -ex
  sudo apt install qperf
  pkill qperf # kill if it is already running
  qperf&
  echo "Running qperf server..."
  echo "You may now continue with the other script."
EOF

