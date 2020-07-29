#!/bin/bash

# start websocket and server
./go/src/github.com/alexanderkarlis/naive-blockchain
status=$?
if [ $status -ne 0 ]; then
  echo "Failed to start go binary: $status"
  exit $status
fi

# Start the second process
npm run start /app/ -D
status=$?
if [ $status -ne 0 ]; then
  echo "Failed to start my_second_process: $status"
  exit $status
fi