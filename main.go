package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	nameFolder := os.Getenv("NAME_FOLDER")
	listenAddr := os.Getenv("LISTEN_ADDR")

	s := NewServer(listenAddr, nameFolder)
	s.Run()
}
