package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	if os.Getenv("ETHERSCAN_API_KEY") == "" {
		log.Fatal("ETHERSCAN_API_KEY is required")
	}
}
