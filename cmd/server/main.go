package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/nikolai/ai-resume-builder/backend/internal/config"
	"github.com/nikolai/ai-resume-builder/backend/internal/database"
	"github.com/nikolai/ai-resume-builder/backend/internal/handlers"
	"github.com/nikolai/ai-resume-builder/backend/internal/repository"
	"github.com/nikolai/ai-resume-builder/backend/internal/router"
	"github.com/nikolai/ai-resume-builder/backend/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Initialize database configuration
	dbConfig := config.NewDatabaseConfig()

	// Create database connection
	db, err := database.NewDB(dbConfig.ConnectionString())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo, db)
	keywordService := service.NewKeywordService(db)
	llmService := service.NewLLMService()
	resumeService := service.NewResumeService(db, keywordService, llmService)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	resumeHandler := handlers.NewResumeHandler(resumeService)

	// Setup router
	r := router.SetupRouter(userHandler, resumeHandler)

	// Start server
	go func() {
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create shutdown context with timeout
	_, cancel := context.WithTimeout(context.Background(), 5)
	defer cancel()

	log.Println("Server exiting")
}
