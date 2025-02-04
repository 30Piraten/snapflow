package config

import (
	"log"

	"github.com/joho/godotenv"
)

func Env() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env file")
	}
}
