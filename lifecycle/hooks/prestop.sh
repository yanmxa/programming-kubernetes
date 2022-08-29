#!/bin/sh

set -e pipefail

echo ""
time=$(date "+%Y-%m-%d %H:%M:%S")
echo "====================================== prestop start: [ $time ] ================================================"
echo "start time: $time" > /usr/share/prestop

second=5
while [ $second -gt 1 ]; do
  sleep 1
  time=$(date "+%Y-%m-%d %H:%M:%S")
  echo "$time prestop is processing..."
  second=`expr $second - 1`
done

time=$(date "+%Y-%m-%d %H:%M:%S")
echo "======================================  prestop end: [ $time ]  ================================================"
echo "end time: $time" >> /usr/share/prestop
