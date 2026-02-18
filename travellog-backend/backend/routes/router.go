package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/travellog/backend/config"
	"github.com/travellog/backend/handlers"
	"github.com/travellog/backend/middleware"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// basic health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// auth
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register(cfg))
		auth.POST("/login", handlers.Login(cfg))
	}

	// protected example (for demo)
	protected := r.Group("/")
	protected.Use(middleware.JWTAuth(cfg))
	{
		protected.GET("/me", handlers.Me())
	}

	return r
}
