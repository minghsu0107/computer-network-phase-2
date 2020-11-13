package main

import (
	"context"
	"log"
	"net"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

func main() {
	ln, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		log.Print("error listening on port 8000")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print("failed to accept connection")
		}
		go handleConnection(conn)
	}
}
