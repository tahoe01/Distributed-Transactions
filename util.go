package main

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Server struct {
	Ip   string
	Port int
}

func Check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func ReadConfigFile(configFile string) map[string]Server {
	branchMap := make(map[string]Server)
	data, err := ioutil.ReadFile(configFile)
	Check(err)

	line := strings.Split(string(data), "\n")
	for _, v := range line {
		branchInfo := strings.Split(v, " ")
		branch, ip, port := branchInfo[0], branchInfo[1], branchInfo[2]
		portNum, err := strconv.Atoi(port)
		Check(err)

		branchMap[branch] = Server{ip, portNum}
	}
	return branchMap
}
