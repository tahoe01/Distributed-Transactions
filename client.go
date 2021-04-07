package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var branchMap map[string]Server

func readCommand() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Printf("Command: %v\n", scanner.Text())
	}
}

func main() {
	fmt.Printf("Client process started.\n")
	log.SetPrefix("Client: ")
	log.SetFlags(0)

	if len(os.Args) != 2 {
		fmt.Println("ERROR: not enough arguments. Usage: ./client [CONFIG_FILE_PATH]")
		return
	}
	configFile := os.Args[1]
	branchMap = ReadConfigFile(configFile)

	readCommand()
}
