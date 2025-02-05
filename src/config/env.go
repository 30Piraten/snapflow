package config

import (
	"log"

	"github.com/joho/godotenv"
)

// Env initializes and loads rhe .env file using the godotenv package
func Env() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}
}
