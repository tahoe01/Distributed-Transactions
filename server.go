package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/dayuebai/Distributed-Transactions/utils"
)

var branchMap map[string]utils.Server
var branchId, configFile string
var port int
var err error

func main() {
	fmt.Printf("Server process started.\n")
	log.SetPrefix("Server: ")
	log.SetFlags(0)

	if len(os.Args) != 4 {
		fmt.Println("ERROR: not enough arguments. Usage: ./server [BRANCH] [PORT] [CONFIG_FILE_PATH]")
		return
	}
	branchId, configFile = os.Args[1], os.Args[3]
	branchMap = utils.ReadConfigFile(configFile)
	port, err = strconv.Atoi(os.Args[2])
	utils.Check(err)
}
