package main

import (
	"fmt"
	"net"
	"log"
)

const status = "OK\n"

func main() {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	var msg []byte
	conn.Read(msg)
	if string(msg) != status {
		fmt.Errorf("Wrong msg: %v", string(msg))
	}
}
