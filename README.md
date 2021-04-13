# Distributed-Transactions

## Overview 
A distributed transaction system that supports transactions which read and write to distributed objects while ensuring full **ACI(D)** properties.

## Project Structure

* `client.go` Code of client process

* `server.go` Code of server process
  
* `util.go` Code of utility function

* `run.sh` Executable script that can run either client or server processes at a time by configuration

* `branch.conf` Configuration file

* `test/` Files under this directory are test files

## Run

* Run client process

  * `./run.sh "client" <client ID> <config file path>`
  * e.g. `./run.sh "client" 1 branch.conf`


* Run Server process

  * `./run.sh "server" <branch ID> <config file path>`
  * e.g. `./run.sh "server" "A" branch.conf`
  
## Contact

Dayue Bai <dayueb2@illinois.edu>


