package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
)

var branchMap map[string]Server
var branchId, configFile string
var port int
var err error

type Handler struct { }

type Request struct {
	Operation int // 1 for DEPOSIT, 2 for BALANCE, 3 for WITHDRAW, 4 for COMMIT, 5 for ABORT
	Account string
	Amount int
}

func (h *Handler) ExecCmd(args *Request, reply *string) error {
	*reply = "RPC works"
	return nil
}

func main() {
	fmt.Printf("Server process started.\n")
	log.SetPrefix("Server: ")
	log.SetFlags(0)

	if len(os.Args) != 4 {
		fmt.Println("ERROR: not enough arguments. Usage: ./server [BRANCH] [PORT] [CONFIG_FILE_PATH]")
		return
	}
	branchId, configFile = os.Args[1], os.Args[3]
	branchMap = ReadConfigFile(configFile)
	port, err = strconv.Atoi(os.Args[2])
	Check(err)

	handler := Handler{}
	rpc.Register(&handler)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	Check(e)
	http.Serve(l, nil)
}
