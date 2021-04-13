#!/bin/bash

if [ $1 == "client" ]; then
  go run client.go util.go $2 $3 # ./client <clientID> <config>
  #cat test/test1 | ./client branch.conf 1
else
  go run server.go util.go $2 $3 # ./server <branch> <config>
fi