package main

import (
	"hse-results/routes"
	"log"

	"hse-results/database"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() // Load .env file
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	err = database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer database.DisconnectDB()
	routes.SetupRouter()
}
