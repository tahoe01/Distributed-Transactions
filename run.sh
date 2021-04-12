#!/bin/bash

if [ $1 == "client" ]; then
  go build client.go util.go
  ./client $2 $3 # ./client <config> <branchID>
  #cat test/test1 | ./client branch.conf 1
else
  go build server.go util.go
  ./server $2 $3 $4 # ./server <branch> <port> <config>
fi