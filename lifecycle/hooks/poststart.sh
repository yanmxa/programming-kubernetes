#!/bin/bash

set -e pipefail

time=$(date "+%Y-%m-%d %H:%M:%S")
echo "===================================== poststart start: [ $time ] ============================================="
echo "start time: $time" > /usr/share/poststart

second=5
while [[ $second -ne 0 ]]; do
  sleep 1
  time=$(date "+%Y-%m-%d %H:%M:%S")
  echo " [ $time ] poststart is processing..."
  ((second--))
done

time=$(date "+%Y-%m-%d %H:%M:%S")
echo "===================================== poststart end: [ $time ]  =============================================="
echo "end time: $time" >> /usr/share/poststart

