package routes

import (
	"go-gin-postgre-crud/controllers"
	"go-gin-postgre-crud/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Apply CORS globally — handles preflight OPTIONS on all routes
	router.Use(middleware.CORSMiddleware())

	// Public health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Service is healthy",
		})
	})

	// ─── Public Auth Routes ───────────────────────────────────────────────────
	auth := router.Group("/api/auth")
	{
		auth.POST("/register", controllers.Register) // POST /api/auth/register
		auth.POST("/login", controllers.Login)       // POST /api/auth/login
	}

	// ─── Protected Routes (JWT required) ─────────────────────────────────────
	api := router.Group("/api", middleware.AuthMiddleware())
	{
		// Auth
		api.GET("/auth/me", controllers.Me) // GET /api/auth/me

		// PRODUCT ROUTES
		api.GET("products", controllers.GetProducts)          // GET    /api/products
		api.GET("products/:id", controllers.GetProductByID)   // GET    /api/products/:id
		api.POST("products", controllers.CreateProduct)       // POST   /api/products
		api.PUT("products/:id", controllers.UpdateProduct)    // PUT    /api/products/:id
		api.DELETE("products/:id", controllers.DeleteProduct) // DELETE /api/products/:id
	}
}
