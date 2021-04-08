#!/bin/bash

go build client.go util.go
cat test/test1 | ./client branch.conf 1
