package controllers

import (
	"net/http"

	"go-gin-postgre-crud/config"
	"go-gin-postgre-crud/models"
	"go-gin-postgre-crud/utils"

	"github.com/gin-gonic/gin"
)

// GetContainers handles the GET request to fetch all containers
// @Summary      List all containers
// @Description  Get a list of all containers, including related master data (type, shipping line, port, voyage, user).
// @Tags         Containers
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Container
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /api/containers [get]
func GetContainers(c *gin.Context) {
	var containers []models.Container
	if err := config.DB.Preload("ContainerType").Preload("ShippingLine").Preload("Port").Preload("Voyage").Preload("User").Find(&containers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, containers)
}

// GetContainerByID handles the GET request to fetch a specific container by ID
// @Summary      Get container by ID
// @Description  Retrieve a single container by its ID, including related master data.
// @Tags         Containers
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Container ID"
// @Success      200 {object} models.Container
// @Failure      404 {object} map[string]string "Container not found"
// @Router       /api/containers/{id} [get]
func GetContainerByID(c *gin.Context) {
	id := c.Param("id")
	var container models.Container
	if err := config.DB.Preload("ContainerType").Preload("ShippingLine").Preload("Port").Preload("Voyage").Preload("User").First(&container, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Container not found"})
		return
	}
	c.JSON(http.StatusOK, container)
}

// CreateContainer handles the POST request to create a new container or multiple containers
// @Summary      Create container(s)
// @Description  Create a single container or a batch of containers. The container number must be valid ISO 6346 and unique.
// @Tags         Containers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object true "Container payload (object or array of objects)"
// @Success      201 {object} map[string]interface{} "Container(s) created successfully"
// @Failure      400 {object} map[string]string "Bad request or invalid container number"
// @Failure      409 {object} map[string]string "Container number already exists"
// @Router       /api/containers [post]
func CreateContainer(c *gin.Context) {
	var rawData interface{}
	if err := c.ShouldBindJSON(&rawData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if it's an array
	if containersArray, ok := rawData.([]interface{}); ok && len(containersArray) > 0 {
		var containers []models.Container
		var containerNumbers []string
		seenNumbers := make(map[string]bool)

		// Convert interface{} array to Container slice and collect numbers
		for _, item := range containersArray {
			if m, ok := item.(map[string]interface{}); ok {
				containerNumber := m["container_number"].(string)

				// Validate ISO 6346 format
				if !utils.IsValidISO6346(containerNumber) {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ISO 6346 container number: " + containerNumber})
					return
				}

				// Check for duplicates within the request array
				if seenNumbers[containerNumber] {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate container number found in request: " + containerNumber})
					return
				}
				seenNumbers[containerNumber] = true
				containerNumbers = append(containerNumbers, containerNumber)

				container := models.Container{
					ContainerNumber: containerNumber,
					ContainerTypeID: uint(m["container_type_id"].(float64)),
					ShippingLineID:  uint(m["shipping_line_id"].(float64)),
					PortID:          uint(m["port_id"].(float64)),
					VoyageID:        uint(m["voyage_id"].(float64)),
					Status:          m["status"].(string),
					CreatedBy:       uint(m["created_by"].(float64)),
				}
				containers = append(containers, container)
			}
		}

		// Check for existing container numbers in database
		var existingCount int64
		if err := config.DB.Model(&models.Container{}).Where("container_number IN ?", containerNumbers).Count(&existingCount).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if existingCount > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "One or more container numbers already exist in the database"})
			return
		}

		if err := config.DB.CreateInBatches(containers, 1000).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "Containers created successfully",
			"count":   len(containers),
			"data":    containers,
		})
		return
	}

	// Single container
	if containerMap, ok := rawData.(map[string]interface{}); ok {
		container := models.Container{
			ContainerNumber: containerMap["container_number"].(string),
			ContainerTypeID: uint(containerMap["container_type_id"].(float64)),
			ShippingLineID:  uint(containerMap["shipping_line_id"].(float64)),
			PortID:          uint(containerMap["port_id"].(float64)),
			VoyageID:        uint(containerMap["voyage_id"].(float64)),
			Status:          containerMap["status"].(string),
			CreatedBy:       uint(containerMap["created_by"].(float64)),
		}

		// Validate ISO 6346 format
		if !utils.IsValidISO6346(container.ContainerNumber) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ISO 6346 container number"})
			return
		}

		// Check for duplicate container number
		var existingContainer models.Container
		if err := config.DB.Where("container_number = ?", container.ContainerNumber).First(&existingContainer).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Container number already exists"})
			return
		}

		if err := config.DB.Create(&container).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, container)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
}

// UpdateContainer handles the PUT request to update an existing container
// @Summary      Update container
// @Description  Update an existing container by ID. Validates ISO 6346 format if container number is provided.
// @Tags         Containers
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Container ID"
// @Param        request body models.Container true "Container payload"
// @Success      200 {object} map[string]string "Container updated successfully"
// @Failure      400 {object} map[string]string "Bad request or invalid container number"
// @Router       /api/containers/{id} [put]
func UpdateContainer(c *gin.Context) {
	id := c.Param("id")
	var container models.Container
	if err := c.ShouldBindJSON(&container); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate ISO 6346 format if container number is provided
	if container.ContainerNumber != "" && !utils.IsValidISO6346(container.ContainerNumber) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ISO 6346 container number"})
		return
	}

	if err := config.DB.Where("id = ?", id).Updates(&container).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Container updated successfully"})
}

// DeleteContainer handles the DELETE request to remove a container
// @Summary      Delete container
// @Description  Delete a container by ID.
// @Tags         Containers
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Container ID"
// @Success      200 {object} map[string]string "Container deleted successfully"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /api/containers/{id} [delete]
func DeleteContainer(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Container{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Container deleted successfully"})
}
