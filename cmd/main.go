package main

import (
	"github.com/joho/godotenv"
	"log"
	"message-service/internal/config"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	_ = config.Read()
	return nil
}
