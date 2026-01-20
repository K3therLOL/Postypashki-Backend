package main

import (
	"cryptoserver/rest"
	"log"
)

func main() {
	log.Println("Server started on port :8080")

	if err := rest.CreateAndRun(); err != nil {
		log.Fatal(err)
	}
}
