#!/bin/bash

if [ $1 == "client" ]; then
  go run client.go util.go $2 $3 # ./client <clientID> <config>
  #  The following commands are used to test by reading files
  #  go build client.go util.go
  #  cat test/test1 | ./client 1 branch.conf
else
  go run server.go util.go $2 $3 # ./server <branch> <config>
fi