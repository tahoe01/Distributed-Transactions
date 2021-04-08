package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
)

var clientID string
var branchMap map[string]Server
var isBegin = false

func pickRandomCoordinator() string {
	servers := [5]string{"A", "B", "C", "D", "E"} // hard coding (being lazy)
	return servers[rand.Intn(5)]
}

func readCommand() {
	var coordinator string
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		switch cmd := scanner.Text(); cmd {
		case "BEGIN":
			isBegin = true
			coordinator = pickRandomCoordinator()
			fmt.Println("OK")
		case "ABORT":
			// TODO: send cmd to coordinator
			isBegin = false
		case "COMMIT":
			// TODO: send cmd to coordinator
			isBegin = false
		default:
			if isBegin {
				// TODO: send cmd to coordinator
				fmt.Printf("Coordinator: %v\n", coordinator)
			}
		}
	}
}

func main() {
	fmt.Printf("Client process started.\n")
	log.SetPrefix("Client: ")
	log.SetFlags(0)

	if len(os.Args) != 3 {
		fmt.Println("ERROR: not enough arguments. Usage: ./client [CONFIG_FILE_PATH]")
		return
	}
	configFile := os.Args[1]
	clientID = os.Args[2]
	branchMap = ReadConfigFile(configFile)

	readCommand()
}
