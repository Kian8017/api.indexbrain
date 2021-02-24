package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"fmt"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisUrl := os.Getenv("REDIS_URL")
	sqlitePath := os.Getenv("SQLITE_PATH")

	fmt.Println("Got redis", redisUrl, "sqlite", sqlitePath);
}
