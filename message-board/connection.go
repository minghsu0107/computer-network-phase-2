package main

import (
	"fmt"
	"log"
	"net"
)

func handleConnection(conn net.Conn) error {
	defer conn.Close()
	requestStr := make([]byte, 3000)
	_, err := conn.Read(requestStr)
	if err != nil {
		log.Print("failed to read request contents")
		return err
	}
	fmt.Println(string(requestStr))
	req := parseRequest(string(requestStr))
	return handle(conn, req)
}
