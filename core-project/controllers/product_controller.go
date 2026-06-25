package controllers

import (
	"net/http"

	"go-gin-postgre-crud/config"
	"go-gin-postgre-crud/models"

	"github.com/gin-gonic/gin"
)

// GetProducts handles the GET request to fetch all products
func GetProducts(c *gin.Context) {
	var products []models.Product
	if err := config.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetProductByID handles the GET request to fetch a specific product by ID
func GetProductByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}

// CreateProduct handles the POST request to create a new product or multiple products
func CreateProduct(c *gin.Context) {
	var rawData interface{}
	if err := c.ShouldBindJSON(&rawData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if it's an array
	if productsArray, ok := rawData.([]interface{}); ok && len(productsArray) > 0 {
		var products []models.Product
		// Convert interface{} array to Product slice
		for _, item := range productsArray {
			if m, ok := item.(map[string]interface{}); ok {
				product := models.Product{
					Name:  m["name"].(string),
					Price: int(m["price"].(float64)),
					Stock: int(m["stock"].(float64)),
				}
				products = append(products, product)
			}
		}

		if err := config.DB.CreateInBatches(products, 1000).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "Products created successfully",
			"count":   len(products),
			"data":    products,
		})
		return
	}

	// Single product
	if productMap, ok := rawData.(map[string]interface{}); ok {
		product := models.Product{
			Name:  productMap["name"].(string),
			Price: int(productMap["price"].(float64)),
			Stock: int(productMap["stock"].(float64)),
		}
		if err := config.DB.Create(&product).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, product)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
}

// UpdateProduct handles the PUT request to update an existing product
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Where("id = ?", id).Updates(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}

// DeleteProduct handles the DELETE request to remove a product
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Product{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}
