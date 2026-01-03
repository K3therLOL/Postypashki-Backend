package main

import (
	"cryptoserver/http"
	"log"
)

func main() {
	log.Println("Server started on port :8080")

	if err := http.CreateAndRun(); err != nil {
		log.Fatal(err)
	}
}
