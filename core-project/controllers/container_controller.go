package controllers

import (
	"net/http"

	"go-gin-postgre-crud/config"
	"go-gin-postgre-crud/models"

	"github.com/gin-gonic/gin"
)

// GetContainers handles the GET request to fetch all containers
func GetContainers(c *gin.Context) {
	var containers []models.Container
	if err := config.DB.Preload("ContainerType").Preload("ShippingLine").Preload("Port").Preload("Voyage").Preload("User").Find(&containers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, containers)
}

// GetContainerByID handles the GET request to fetch a specific container by ID
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
func CreateContainer(c *gin.Context) {
	var rawData interface{}
	if err := c.ShouldBindJSON(&rawData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if it's an array
	if containersArray, ok := rawData.([]interface{}); ok && len(containersArray) > 0 {
		var containers []models.Container
		// Convert interface{} array to Container slice
		for _, item := range containersArray {
			if m, ok := item.(map[string]interface{}); ok {
				container := models.Container{
					ContainerNumber: m["container_number"].(string),
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
func UpdateContainer(c *gin.Context) {
	id := c.Param("id")
	var container models.Container
	if err := c.ShouldBindJSON(&container); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Where("id = ?", id).Updates(&container).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Container updated successfully"})
}

// DeleteContainer handles the DELETE request to remove a container
func DeleteContainer(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Container{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Container deleted successfully"})
}
