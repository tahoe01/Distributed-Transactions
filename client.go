package main

import (
	"fmt"
	"os"
)

type Server struct {
	Ip string
	Port int
}

var branchMap map[string]Server

func readConfigFile(configFile string) {
	fmt.Println(configFile)
}

func main() {
	fmt.Printf("Client process started.\n")

	if len(os.Args) != 2 {
		fmt.Println("ERROR: not enough arguments. Usage: ./client [CONFIG_FILE_PATH]")
		return
	}
	configFile := os.Args[1]
	readConfigFile(configFile)
}
