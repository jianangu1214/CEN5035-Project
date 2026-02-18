package main

import (
	"log"

	"github.com/travellog/backend/config"
	"github.com/travellog/backend/database"
	"github.com/travellog/backend/routes"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup router
	router := routes.SetupRouter(cfg)

	// Start server
	log.Printf("Server starting on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
