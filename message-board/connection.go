package main

import (
	"fmt"
	"log"
	"net"
)

var requestLen int

func handleConnection(conn net.Conn) error {
	var err error
	defer conn.Close()
	requestStr := make([]byte, 30000)
	requestLen, err = conn.Read(requestStr)
	if err != nil {
		log.Print("failed to read request contents")
		return err
	}
	fmt.Println(string(requestStr))
	req := parseRequest(string(requestStr))
	return handle(conn, req)
}
