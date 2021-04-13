package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"time"
	//"math/rand"
	"net/rpc"
	"os"
	"strconv"
	"strings"
)

var clientID string
var branchMap map[string]Server

type Request struct {
	TransactionID string
	Operation int // 1 for DEPOSIT, 2 for BALANCE, 3 for WITHDRAW, 4 for COMMIT, 5 for ABORT
	Account string
	Amount int
}

type Reply struct {
	Status int // 0 for FAIL, 1 for SUCCESS
	Msg string
}

func generateTransactionID() string {
	nanoSecond := time.Now().UnixNano()
	return fmt.Sprintf("%v.%v", strconv.FormatInt(nanoSecond, 10), clientID) // transaction ID Format: TS.clientID
}

func pickRandomCoordinator() string {
	servers := [5]string{"A", "B", "C", "D", "E"} // hard coding (being lazy)
	return servers[rand.Intn(5)]
}

func parseCmd(cmd string, transactionID string) Request {
	if components := strings.Split(cmd, " "); len(components) == 2 {
		return Request{transactionID, 2,components[1], -1 } // For BALANCE cmd, amount set to -1
	} else {
		amount, err := strconv.Atoi(components[2])
		Check(err)
		if components[0] == "DEPOSIT" {
			return Request{transactionID, 1, components[1], amount}
		} else { // WITHDRAW
			return Request{transactionID, 3, components[1], amount * (-1)}
		}
	}
}

func readCommand() {
	var coordinator, transactionID string
	var client *rpc.Client
	var err error
	var accountWrite []string
	isBegin := false
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		var req Request
		var reply Reply
		switch cmd := scanner.Text(); cmd {
		case "BEGIN":
			if !isBegin {
				isBegin = true
				accountWrite = nil
				transactionID = generateTransactionID()
				coordinator = pickRandomCoordinator()
				client, err = rpc.DialHTTP("tcp", fmt.Sprintf("%v:%v", branchMap[coordinator].Ip, branchMap[coordinator].Port))
				Check(err)
				fmt.Println("OK")
			}
		case "ABORT":
			if isBegin {
				req = Request{transactionID, 5, "", -1}
				err = client.Call("Handler.DeliverCmd", &req, &reply)
				Check(err)
				isBegin = false
				fmt.Println(reply.Msg)
			}
		case "COMMIT":
			if isBegin {
				req = Request{transactionID, 4, strings.Join(accountWrite, ","), -1}
				err = client.Call("Handler.DeliverCmd", &req, &reply)
				Check(err)
				isBegin = false
				fmt.Println(reply.Msg)
			}
		default:
			if isBegin {
				info := strings.Split(cmd, " ")
				if cmdType := info[0]; cmdType == "DEPOSIT" || cmdType == "WITHDRAW" { // Write operations
					accountWrite = append(accountWrite, info[1])
				}
				req = parseCmd(cmd, transactionID)
				err = client.Call("Handler.DeliverCmd", &req, &reply)
				Check(err)
				if reply.Status == -1 {
					isBegin = false
				}
				fmt.Println(reply.Msg)
			}
		}
	}
}

func main() {
	fmt.Printf("Client process started.\n")
	log.SetPrefix("Client: ")
	log.SetFlags(0)

	if len(os.Args) != 3 {
		fmt.Println("ERROR: not enough arguments. Usage: ./client [ClientID] [CONFIG_FILE_PATH]")
		return
	}
	clientID = os.Args[1]
	configFile := os.Args[2]
	branchMap = ReadConfigFile(configFile)

	readCommand()
}
