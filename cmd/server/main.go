package main

import (
	"log"

	"github.com/nikolai/ai-resume-builder/backend/config"
	"github.com/nikolai/ai-resume-builder/backend/internal/handlers"
	"github.com/nikolai/ai-resume-builder/backend/internal/utils"
)

func main() {
	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := utils.NewLogger()

	// Initialize router
	router := handlers.SetupRouter(logger)

	// Start server
	logger.Infof("Starting server on port %s", cfg.ServerPort)
	if err := router.Run(":" + cfg.ServerPort); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
