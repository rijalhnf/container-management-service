package controllers

import (
	"net/http"

	"go-gin-postgre-crud/config"
	"go-gin-postgre-crud/models"

	"github.com/gin-gonic/gin"
)

// GetContainerTypes handles the GET request to fetch all container types
// @Summary      List all container types
// @Description  Get a list of all container types.
// @Tags         Container Types
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.ContainerType
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /api/container-types [get]
func GetContainerTypes(c *gin.Context) {
	var containerTypes []models.ContainerType
	if err := config.DB.Find(&containerTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, containerTypes)
}

// GetContainerTypeByID handles the GET request to fetch a specific container type by ID
// @Summary      Get container type by ID
// @Description  Retrieve a single container type by its ID.
// @Tags         Container Types
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Container type ID"
// @Success      200 {object} models.ContainerType
// @Failure      404 {object} map[string]string "Container type not found"
// @Router       /api/container-types/{id} [get]
func GetContainerTypeByID(c *gin.Context) {
	id := c.Param("id")
	var containerType models.ContainerType
	if err := config.DB.First(&containerType, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Container type not found"})
		return
	}
	c.JSON(http.StatusOK, containerType)
}

// CreateContainerType handles the POST request to create a new container type
// @Summary      Create container type
// @Description  Create a new container type. The `code` field must be unique.
// @Tags         Container Types
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.ContainerType true "Container type payload"
// @Success      201 {object} models.ContainerType
// @Failure      400 {object} map[string]string "Bad request"
// @Failure      409 {object} map[string]string "Container type code already exists"
// @Router       /api/container-types [post]
func CreateContainerType(c *gin.Context) {
	var containerType models.ContainerType
	if err := c.ShouldBindJSON(&containerType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate code
	var existing models.ContainerType
	if err := config.DB.Where("code = ?", containerType.Code).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Container type code already exists"})
		return
	}

	if err := config.DB.Create(&containerType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, containerType)
}

// UpdateContainerType handles the PUT request to update an existing container type
// @Summary      Update container type
// @Description  Update an existing container type by ID.
// @Tags         Container Types
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Container type ID"
// @Param        request body models.ContainerType true "Container type payload"
// @Success      200 {object} map[string]string "Container type updated successfully"
// @Failure      400 {object} map[string]string "Bad request"
// @Failure      409 {object} map[string]string "Container type code already exists"
// @Router       /api/container-types/{id} [put]
func UpdateContainerType(c *gin.Context) {
	id := c.Param("id")
	var containerType models.ContainerType
	if err := c.ShouldBindJSON(&containerType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate code (excluding current ID)
	var existing models.ContainerType
	if err := config.DB.Where("code = ? AND id <> ?", containerType.Code, id).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Container type code already exists"})
		return
	}

	if err := config.DB.Where("id = ?", id).Updates(&containerType).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Container type updated successfully"})
}

// DeleteContainerType handles the DELETE request to remove a container type
// @Summary      Delete container type
// @Description  Delete a container type by ID.
// @Tags         Container Types
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Container type ID"
// @Success      200 {object} map[string]string "Container type deleted successfully"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /api/container-types/{id} [delete]
func DeleteContainerType(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.ContainerType{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Container type deleted successfully"})
}
