// @title           Container Management API
// @version         1.0
// @description     REST API untuk manajemen container dengan autentikasi JWT.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@container.local

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8888
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"go-gin-postgre-crud/config"
	"go-gin-postgre-crud/models"
	"go-gin-postgre-crud/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.User{}, 
		&models.ContainerType{}, &models.ShippingLine{}, &models.Port{}, &models.Voyage{}, &models.Container{})

	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	routes.SetupRoutes(r)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8090"
	}
	fmt.Printf("Server is running on port %s\n", port)

	// Configure HTTP server with timeout settings
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  10 * time.Second, // Time to read entire request
		WriteTimeout: 10 * time.Second, // Time to write entire response
		IdleTimeout:  60 * time.Second, // Time for keep-alive connection
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
