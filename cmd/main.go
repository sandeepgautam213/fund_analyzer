package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sandeepgautam213/eth-flow-analyzer/config"
	"github.com/sandeepgautam213/eth-flow-analyzer/internal/handler"
)

func main() {
	// Load environment variables from .env file
	config.LoadEnv()

	// Create a new Gin router
	router := gin.Default()

	// Register the /beneficiary endpoint
	router.GET("/beneficiary", handler.GetBeneficiaries)
	router.GET("/payer", handler.GetPayers)

	// Get PORT from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on port", port)
	router.Run(":" + port)
}
