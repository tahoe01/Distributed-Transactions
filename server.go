package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strings"
	"sync"
)

var branchMap map[string]Server
var branchId, configFile string
var port string
var accountMap = AccountMap{&sync.Mutex{}, make(map[string]Account)}

type Handler struct { } // RPC handler

type void struct { } // empty struct

type Request struct {
	TransactionID string
	Operation int // 1 for DEPOSIT, 2 for BALANCE, 3 for WITHDRAW, 4 for COMMIT, 5 for ABORT
	Account string
	Amount int
}

type Reply struct {
	Status int // 0 for SUCCESS, -1 for FAILURE
	Msg string
}

type Account struct {
	mu *sync.Mutex
	accountName string // format: A.xyz
	committedValue int
	transactionID string
	rtsList map[string]void // List of transaction ids (timestamps) that have read the committed value.
	twMap map[string]int // List of tentative writes by the corresponding transaction ids.
}

type AccountMap struct {
	mu *sync.Mutex
	aMap map[string]Account
}


// Local
func getMaxRTs(rtsList map[string]void) string {
	var maxRTs = ""
	for rts := range rtsList {
		if rts > maxRTs {
			maxRTs = rts
		}
	}
	return maxRTs
}

func getMaxWTs(twMap map[string]int, base string, threshold string) string {
	maxWTS := base
	for wts := range twMap {
		if wts > maxWTS && wts <= threshold {
			maxWTS = wts
		}
	}
	return maxWTS
}

func readAccount(args *Request, reply *Reply) {
	accountName := args.Account
	accountMap.mu.Lock()
	account, isExist := accountMap.aMap[accountName]
	accountMap.mu.Unlock()

	if !isExist { // Account not found
		*reply = Reply{-1, "NOT FOUND"}
	} else { // Account exists
		fmt.Println("account exists...")
		account.mu.Lock()
		fmt.Println("Pass lock")
		accountTID, curTID := account.transactionID, args.TransactionID
		maxWts := getMaxWTs(account.twMap, account.transactionID, curTID)
		account.mu.Unlock()
		if curTID > accountTID {
			if maxWts == accountTID { // version of object with maxWTs is committed => read this version and add curTID to rtsList
				account.mu.Lock()
				account.rtsList[curTID] = void{}
				account.mu.Unlock()
				*reply = Reply{0, fmt.Sprintf("%v = %v", args.Account, account.committedValue)}
			} else if maxWts == curTID { // Write by curTID => read version of itself
				*reply = Reply{0, fmt.Sprintf("%v = %v", args.Account, account.committedValue + account.twMap[curTID])}
			} else { // TODO: wait until another transaction is committed or aborted and reapply the read rule
				fmt.Println("case3...")
				fmt.Println("Waiting for another transaction")
			}
		} else {
			*reply = Reply{-1, "ABORTED"}
		}
	}
}

func writeAccount(args *Request, reply *Reply) {
	accountName := args.Account
	accountMap.mu.Lock()
	account, isExist := accountMap.aMap[accountName]
	accountMap.mu.Unlock()
	if !isExist {
		if args.Operation == 1 { // deposit
			newAccount := Account{&sync.Mutex{}, accountName, 0, "", make(map[string]void), make(map[string]int)}
			newAccount.twMap[args.TransactionID] = args.Amount
			accountMap.mu.Lock()
			accountMap.aMap[accountName] = newAccount
			accountMap.mu.Unlock()
			*reply = Reply{0, "OK"}
		} else { // withdraw
			*reply = Reply{-1, "NOT FOUND"}
		}
	} else { // account exists, perform write
		account.mu.Lock()
		curTID, maxRts := args.TransactionID, getMaxRTs(account.rtsList)
		if curTID >= maxRts && curTID > account.transactionID {
			if change, ok := account.twMap[curTID]; ok { // curTID already has an entry in TW list
				account.twMap[curTID] = change + args.Amount
			} else {
				account.twMap[curTID] = args.Amount
			}
			*reply = Reply{0, "OK"}
		} else {
			*reply = Reply{-1, "ABORTED"}
		}
		account.mu.Unlock()
	}
}

// RPC
func (h *Handler) RunCmd(args *Request, reply *Reply) error {
	fmt.Printf("[RUN CMD] TransactionID: %v, Operation: %v, Account: %v, Amount: %v\n", args.TransactionID, args.Operation, args.Account, args.Amount)
	if args.Operation == 2 { // BALANCE
		readAccount(args, reply)
	} else { // DEPOSIT/WITHDRAW : write or create
		writeAccount(args, reply)
	}
	return nil
}

func (h *Handler) CanCommit(args *Request, reply *Reply) error { // TODO: Commit a transaction
	return nil
}

func (h *Handler) DoCommit(args *Request, reply *Reply) error { // TODO: Commit a transaction
	return nil
}

func (h *Handler) DoAbort(args *Request, reply *Reply) error { // Abort a transaction
	fmt.Printf("[ABORT] TransactionID: %v, Operation: %v, Account: %v, Amount: %v\n", args.TransactionID, args.Operation, args.Account, args.Amount)
	abortTID := args.TransactionID
	accountMap.mu.Lock()
	defer accountMap.mu.Unlock()
	for name := range accountMap.aMap {
		//account.mu.Lock()
		accountMap.aMap[name].mu.Lock()
		delete(accountMap.aMap[name].rtsList, abortTID) // remove corresponding read from RTS list if exists
		delete(accountMap.aMap[name].twMap, abortTID) // remove corresponding write from TW list if exists
		if accountMap.aMap[name].transactionID == "" && len(accountMap.aMap[name].twMap) == 0 { // Account not committed by any transaction && No wts in TW list => account should be removed
			delete(accountMap.aMap, name)
		}
		//account.mu.Unlock()
		accountMap.aMap[name].mu.Unlock()
	}
	*reply = Reply{-1, "ABORTED"}
	return nil
}

func (h *Handler) DeliverCmd(args *Request, reply *Reply) error {
	fmt.Printf("[DELIVER CMD] TransactionID: %v, Operation: %v, Account: %v, Amount: %v\n", args.TransactionID, args.Operation, args.Account, args.Amount)
	var client *rpc.Client
	var err error
	var branchReply Reply

	if args.Operation <= 3 { // DEPOSIT, WITHDRAW, BALANCE
		branch := strings.Split(args.Account, ".")[0]
		client, _ = rpc.DialHTTP("tcp",fmt.Sprintf("%v:%v", branchMap[branch].Ip, branchMap[branch].Port))
		err = client.Call("Handler.RunCmd", args, &branchReply)
		Check(err)

		fmt.Printf("branch reply: status (%v), msg(%v)\n", branchReply.Status, branchReply.Msg)
		if branchReply.Status == -1 {
			args = &Request{args.TransactionID, 5, "", -1}
		}
	}
	if args.Operation == 4 { // Commit (2PC)
		canCommit := true
		for _, server := range branchMap { // Phase 1: Ask each branch if they can commit the given transaction
			client, _ = rpc.DialHTTP("tcp",fmt.Sprintf("%v:%v", server.Ip, server.Port))
			err = client.Call("Handler.CanCommit", args, &branchReply)
			Check(err)

			fmt.Printf("CanCommit response: %v\n", branchReply.Msg)
			if branchReply.Status == -1 { // There exists one branch that wants to abort the transaction => break out of the loop & abort!
				canCommit = false
				args = &Request{args.TransactionID, 5, "", -1}
				break
			}
		}
		if canCommit {
			for _, server := range branchMap { // Phase 2: Ask each branch to commit the given transaction
				client, _ = rpc.DialHTTP("tcp",fmt.Sprintf("%v:%v", server.Ip, server.Port))
				err = client.Call("Handler.DoCommit", args, &branchReply)
				Check(err)

				fmt.Printf("DoCommit response: %v\n", branchReply.Msg)
			}
		}
	}
	if args.Operation == 5 { // ABORT
		for _, server := range branchMap {
			client, _ = rpc.DialHTTP("tcp",fmt.Sprintf("%v:%v", server.Ip, server.Port))
			err = client.Call("Handler.DoAbort", args, &branchReply)
			Check(err)

			fmt.Printf("abort response: %v\n", branchReply.Msg)
		}
	}

	*reply = Reply{branchReply.Status, branchReply.Msg}
	return nil
}

func server(port string) {
	handler := Handler{}

	rpc.Register(&handler)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":" + port)
	Check(e)
	defer l.Close()

	err := http.Serve(l, nil)
	Check(err)
}

func main() {
	fmt.Printf("Server process started.\n")
	log.SetPrefix("Server: ")
	log.SetFlags(0)

	if len(os.Args) != 4 {
		fmt.Println("ERROR: not enough arguments. Usage: ./server [BRANCH] [PORT] [CONFIG_FILE_PATH]")
		return
	}
	branchId, port, configFile = os.Args[1], os.Args[2], os.Args[3]
	branchMap = ReadConfigFile(configFile)

	server(port)
}
