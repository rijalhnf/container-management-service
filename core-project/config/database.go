package config

import (
	"fmt"
	"log"
	"os"

	"go-gin-postgre-crud/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	DB = database
	fmt.Println("Database connection established")

	// Auto-migrate all models
	err = database.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.ContainerType{},
		&models.ShippingLine{},
		&models.Port{},
		&models.Voyage{},
		&models.Container{},
	)
	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
	fmt.Println("Database migrations completed")
}

// CREATE TABLE users (
// 	id BIGINT AUTO_INCREMENT PRIMARY KEY,
// 	name VARCHAR(255) NOT NULL,
// 	email VARCHAR(255) NOT NULL UNIQUE,
// 	phone VARCHAR(20),
// 	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
// )
