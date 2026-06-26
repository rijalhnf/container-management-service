package controllers

import (
	"net/http"

	"go-gin-postgre-crud/config"
	"go-gin-postgre-crud/models"

	"github.com/gin-gonic/gin"
)

// GetPorts handles the GET request to fetch all ports
// @Summary      List all ports
// @Description  Get a list of all ports.
// @Tags         Ports
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Port
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /api/ports [get]
func GetPorts(c *gin.Context) {
	var ports []models.Port
	if err := config.DB.Find(&ports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ports)
}

// GetPortByID handles the GET request to fetch a specific port by ID
// @Summary      Get port by ID
// @Description  Retrieve a single port by its ID.
// @Tags         Ports
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Port ID"
// @Success      200 {object} models.Port
// @Failure      404 {object} map[string]string "Port not found"
// @Router       /api/ports/{id} [get]
func GetPortByID(c *gin.Context) {
	id := c.Param("id")
	var port models.Port
	if err := config.DB.First(&port, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Port not found"})
		return
	}
	c.JSON(http.StatusOK, port)
}

// CreatePort handles the POST request to create a new port
// @Summary      Create port
// @Description  Create a new port. The `code` field must be unique.
// @Tags         Ports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.Port true "Port payload"
// @Success      201 {object} models.Port
// @Failure      400 {object} map[string]string "Bad request"
// @Failure      409 {object} map[string]string "Port code already exists"
// @Router       /api/ports [post]
func CreatePort(c *gin.Context) {
	var port models.Port
	if err := c.ShouldBindJSON(&port); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate code
	var existing models.Port
	if err := config.DB.Where("code = ?", port.Code).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Port code already exists"})
		return
	}

	if err := config.DB.Create(&port).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, port)
}

// UpdatePort handles the PUT request to update an existing port
// @Summary      Update port
// @Description  Update an existing port by ID.
// @Tags         Ports
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Port ID"
// @Param        request body models.Port true "Port payload"
// @Success      200 {object} map[string]string "Port updated successfully"
// @Failure      400 {object} map[string]string "Bad request"
// @Failure      409 {object} map[string]string "Port code already exists"
// @Router       /api/ports/{id} [put]
func UpdatePort(c *gin.Context) {
	id := c.Param("id")
	var port models.Port
	if err := c.ShouldBindJSON(&port); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate code (excluding current ID)
	var existing models.Port
	if err := config.DB.Where("code = ? AND id <> ?", port.Code, id).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Port code already exists"})
		return
	}

	if err := config.DB.Where("id = ?", id).Updates(&port).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Port updated successfully"})
}

// DeletePort handles the DELETE request to remove a port
// @Summary      Delete port
// @Description  Delete a port by ID.
// @Tags         Ports
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Port ID"
// @Success      200 {object} map[string]string "Port deleted successfully"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /api/ports/{id} [delete]
func DeletePort(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Port{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Port deleted successfully"})
}
