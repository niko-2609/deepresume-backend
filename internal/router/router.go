package router

import (
	"github.com/gin-gonic/gin"
	"github.com/nikolai/ai-resume-builder/backend/internal/handlers"
	"github.com/nikolai/ai-resume-builder/backend/internal/middleware"
)

// SetupRouter configures all the routes for our application
func SetupRouter(userHandler *handlers.UserHandler, resumeHandler *handlers.ResumeHandler) *gin.Engine {
	router := gin.Default()

	// Middleware
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.GET("/email/:email", userHandler.GetUserByEmail)
			users.PUT("/:id", userHandler.UpdateUser)
		}

		// Resume generation route
		v1.POST("/generate", resumeHandler.GenerateResume)

		// Onboarding route
		v1.POST("/onboarding", userHandler.HandleOnboarding)
	}

	return router
}
