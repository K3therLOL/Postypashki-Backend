package main

import (
	"net"
	"log"
	"fmt"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Errorf("Error: %v", err.Error())
		return
	}
	defer listener.Close()

	const msg = "OK\n"
	for {
		conn, err := listener.Accept()	
		if err != nil {
			log.Fatal(err)
		}

		go func(conn net.Conn) {
			conn.Write([]byte(msg))
			conn.Close()
		}(conn)
	}
}
