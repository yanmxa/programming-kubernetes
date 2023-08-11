#!/bin/bash

set -e pipefail

echo ""
time=$(date "+%Y-%m-%d %H:%M:%S")
echo "====================================== prestop start: [ $time ] ================================================"
echo "start time: $time" > /usr/share/prestop

second=10
while [[ $second -ne 0 ]]; do
  sleep 1
  time=$(date "+%Y-%m-%d %H:%M:%S")
  echo " [ $time ] prestop is processing..."
  ((second--))
done

time=$(date "+%Y-%m-%d %H:%M:%S")
echo "======================================  prestop end: [ $time ]  ================================================"
echo "end time: $time" >> /usr/share/prestop

