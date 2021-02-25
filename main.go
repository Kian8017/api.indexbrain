package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisUrl := os.Getenv("REDIS_URL")
	sqlitePath := os.Getenv("SQLITE_PATH")
	listenAddr := os.Getenv("LISTEN_ADDR")


}
