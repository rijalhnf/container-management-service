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
	config.DB.AutoMigrate(&models.User{}, &models.Product{}, &models.ContainerType{}, &models.ShippingLine{}, &models.Port{}, &models.Voyage{}, &models.Container{})

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
