package main

import (
	"github.com/joho/godotenv"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// DB setup
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	// DB init
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":3000")
	server.Run()
}
