package routes

import (
	"go-gin-postgre-crud/controllers"
	"go-gin-postgre-crud/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "go-gin-postgre-crud/docs"
)

func SetupRoutes(router *gin.Engine) {
	// Apply CORS globally — handles preflight OPTIONS on all routes
	router.Use(middleware.CORSMiddleware())

	// Swagger UI — publicly accessible
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

		// Container Types
		containerType := api.Group("/container-types")
		{
			containerType.GET("", controllers.GetContainerTypes)          // GET    /api/container-types
			containerType.GET("/:id", controllers.GetContainerTypeByID)   // GET    /api/container-types/:id
			containerType.POST("", controllers.CreateContainerType)       // POST   /api/container-types
			containerType.PUT("/:id", controllers.UpdateContainerType)    // PUT    /api/container-types/:id
			containerType.DELETE("/:id", controllers.DeleteContainerType) // DELETE /api/container-types/:id
		}

		// Ports
		port := api.Group("/ports")
		{
			port.GET("", controllers.GetPorts)          // GET    /api/ports
			port.GET("/:id", controllers.GetPortByID)   // GET    /api/ports/:id
			port.POST("", controllers.CreatePort)       // POST   /api/ports
			port.PUT("/:id", controllers.UpdatePort)    // PUT    /api/ports/:id
			port.DELETE("/:id", controllers.DeletePort) // DELETE /api/ports/:id
		}

		// Shipping Lines
		shippingLine := api.Group("/shipping-lines")
		{
			shippingLine.GET("", controllers.GetShippingLines)          // GET    /api/shipping-lines
			shippingLine.GET("/:id", controllers.GetShippingLineByID)   // GET    /api/shipping-lines/:id
			shippingLine.POST("", controllers.CreateShippingLine)       // POST   /api/shipping-lines
			shippingLine.PUT("/:id", controllers.UpdateShippingLine)    // PUT    /api/shipping-lines/:id
			shippingLine.DELETE("/:id", controllers.DeleteShippingLine) // DELETE /api/shipping-lines/:id
		}

		// Voyages
		voyage := api.Group("/voyages")
		{
			voyage.GET("", controllers.GetVoyages)          // GET    /api/voyages
			voyage.GET("/:id", controllers.GetVoyageByID)   // GET    /api/voyages/:id
			voyage.POST("", controllers.CreateVoyage)       // POST   /api/voyages
			voyage.PUT("/:id", controllers.UpdateVoyage)    // PUT    /api/voyages/:id
			voyage.DELETE("/:id", controllers.DeleteVoyage) // DELETE /api/voyages/:id
		}

		container := api.Group("/containers")
		{
			container.GET("", controllers.GetContainers)          // GET    /api/containers
			container.GET("/:id", controllers.GetContainerByID)   // GET    /api/containers/:id
			container.POST("", controllers.CreateContainer)       // POST   /api/containers
			container.PUT("/:id", controllers.UpdateContainer)    // PUT    /api/containers/:id
			container.DELETE("/:id", controllers.DeleteContainer) // DELETE /api/containers/:id
		}
	}
}
