#!/bin/bash

go build client.go
cat test/test1 | ./client branch.conf
