package controllers

import (
	"net/http"

	"go-gin-postgre-crud/config"
	"go-gin-postgre-crud/models"

	"github.com/gin-gonic/gin"
)

// GetShippingLines handles the GET request to fetch all shipping lines
// @Summary      List all shipping lines
// @Description  Get a list of all shipping lines.
// @Tags         Shipping Lines
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.ShippingLine
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /api/shipping-lines [get]
func GetShippingLines(c *gin.Context) {
	var shippingLines []models.ShippingLine
	if err := config.DB.Find(&shippingLines).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, shippingLines)
}

// GetShippingLineByID handles the GET request to fetch a specific shipping line by ID
// @Summary      Get shipping line by ID
// @Description  Retrieve a single shipping line by its ID.
// @Tags         Shipping Lines
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Shipping line ID"
// @Success      200 {object} models.ShippingLine
// @Failure      404 {object} map[string]string "Shipping line not found"
// @Router       /api/shipping-lines/{id} [get]
func GetShippingLineByID(c *gin.Context) {
	id := c.Param("id")
	var shippingLine models.ShippingLine
	if err := config.DB.First(&shippingLine, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Shipping line not found"})
		return
	}
	c.JSON(http.StatusOK, shippingLine)
}

// CreateShippingLine handles the POST request to create a new shipping line
// @Summary      Create shipping line
// @Description  Create a new shipping line. The `code` field must be unique.
// @Tags         Shipping Lines
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.ShippingLine true "Shipping line payload"
// @Success      201 {object} models.ShippingLine
// @Failure      400 {object} map[string]string "Bad request"
// @Failure      409 {object} map[string]string "Shipping line code already exists"
// @Router       /api/shipping-lines [post]
func CreateShippingLine(c *gin.Context) {
	var shippingLine models.ShippingLine
	if err := c.ShouldBindJSON(&shippingLine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate code
	var existing models.ShippingLine
	if err := config.DB.Where("code = ?", shippingLine.Code).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Shipping line code already exists"})
		return
	}

	if err := config.DB.Create(&shippingLine).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, shippingLine)
}

// UpdateShippingLine handles the PUT request to update an existing shipping line
// @Summary      Update shipping line
// @Description  Update an existing shipping line by ID.
// @Tags         Shipping Lines
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Shipping line ID"
// @Param        request body models.ShippingLine true "Shipping line payload"
// @Success      200 {object} map[string]string "Shipping line updated successfully"
// @Failure      400 {object} map[string]string "Bad request"
// @Failure      409 {object} map[string]string "Shipping line code already exists"
// @Router       /api/shipping-lines/{id} [put]
func UpdateShippingLine(c *gin.Context) {
	id := c.Param("id")
	var shippingLine models.ShippingLine
	if err := c.ShouldBindJSON(&shippingLine); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate code (excluding current ID)
	var existing models.ShippingLine
	if err := config.DB.Where("code = ? AND id <> ?", shippingLine.Code, id).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Shipping line code already exists"})
		return
	}

	if err := config.DB.Where("id = ?", id).Updates(&shippingLine).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Shipping line updated successfully"})
}

// DeleteShippingLine handles the DELETE request to remove a shipping line
// @Summary      Delete shipping line
// @Description  Delete a shipping line by ID.
// @Tags         Shipping Lines
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Shipping line ID"
// @Success      200 {object} map[string]string "Shipping line deleted successfully"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /api/shipping-lines/{id} [delete]
func DeleteShippingLine(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.ShippingLine{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Shipping line deleted successfully"})
}
