// package routes

// import (
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/travellog/backend/config"
// 	"github.com/travellog/backend/handlers"
// 	"github.com/travellog/backend/middleware"
// )

// func SetupRouter(cfg *config.Config) *gin.Engine {
// 	r := gin.Default()

// 	// basic health check
// 	r.GET("/health", func(c *gin.Context) {
// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
// 	})

// 	// 1. Auth Routes (User Handler)
// 	// 队友负责的部分
// 	auth := r.Group("/auth")
// 	{
// 		auth.POST("/register", handlers.Register(cfg))
// 		auth.POST("/login", handlers.Login(cfg))
// 	}

// 	// 2. Protected Routes
// 	protected := r.Group("/")
// 	protected.Use(middleware.JWTAuth(cfg))
// 	{
// 		protected.GET("/me", handlers.Me())

// 		// --- Hotel Routes (你的部分) ---

// 		hotelHandler := handlers.NewHotelHandler()

// 		hotelGroup := protected.Group("/hotels")
// 		{
// 			hotelGroup.GET("", hotelHandler.GetAllHotels)       // GET /hotels
// 			hotelGroup.POST("", hotelHandler.CreateHotel)       // POST /hotels
// 			hotelGroup.GET("/:id", hotelHandler.GetHotel)       // GET /hotels/:id
// 			hotelGroup.PUT("/:id", hotelHandler.UpdateHotel)    // PUT /hotels/:id
// 			hotelGroup.DELETE("/:id", hotelHandler.DeleteHotel) // DELETE /hotels/:id
// 		}
// 	}

// 	return r
// }

package routes

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/travellog/backend/config"
	"github.com/travellog/backend/handlers"
	"github.com/travellog/backend/middleware"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// CORS: allow frontend (Vite) to call backend APIs
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// basic health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// 1. Auth Routes (User Handler)
	// 队友负责的部分
	auth := r.Group("/auth")
	{
		auth.POST("/register", handlers.Register(cfg))
		auth.POST("/login", handlers.Login(cfg))
	}

	// 2. Protected Routes
	protected := r.Group("/")
	protected.Use(middleware.JWTAuth(cfg))
	{
		protected.GET("/me", handlers.Me())

		// --- Hotel Routes (你的部分) ---
		hotelHandler := handlers.NewHotelHandler()

		hotelGroup := protected.Group("/hotels")
		{
			hotelGroup.GET("", hotelHandler.GetAllHotels)       // GET /hotels
			hotelGroup.POST("", hotelHandler.CreateHotel)       // POST /hotels
			hotelGroup.GET("/:id", hotelHandler.GetHotel)       // GET /hotels/:id
			hotelGroup.PUT("/:id", hotelHandler.UpdateHotel)    // PUT /hotels/:id
			hotelGroup.DELETE("/:id", hotelHandler.DeleteHotel) // DELETE /hotels/:id
		}
	}

	return r
}
