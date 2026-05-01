package main

import (
	"log"
	"taskserver/server"
)

func main() {
	log.Println("Started server at :8080.")
	if err := server.CreateAndRun(); err != nil {
		log.Println("Couldn't run server.")
	}
}
