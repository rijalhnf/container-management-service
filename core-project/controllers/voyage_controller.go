package controllers

import (
	"net/http"

	"go-gin-postgre-crud/config"
	"go-gin-postgre-crud/models"

	"github.com/gin-gonic/gin"
)

// GetVoyages handles the GET request to fetch all voyages
// @Summary      List all voyages
// @Description  Get a list of all voyages, including origin and destination ports.
// @Tags         Voyages
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} models.Voyage
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /api/voyages [get]
func GetVoyages(c *gin.Context) {
	var voyages []models.Voyage
	if err := config.DB.Preload("OriginPort").Preload("DestinationPort").Find(&voyages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, voyages)
}

// GetVoyageByID handles the GET request to fetch a specific voyage by ID
// @Summary      Get voyage by ID
// @Description  Retrieve a single voyage by its ID, including origin and destination ports.
// @Tags         Voyages
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Voyage ID"
// @Success      200 {object} models.Voyage
// @Failure      404 {object} map[string]string "Voyage not found"
// @Router       /api/voyages/{id} [get]
func GetVoyageByID(c *gin.Context) {
	id := c.Param("id")
	var voyage models.Voyage
	if err := config.DB.Preload("OriginPort").Preload("DestinationPort").First(&voyage, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Voyage not found"})
		return
	}
	c.JSON(http.StatusOK, voyage)
}

// CreateVoyage handles the POST request to create a new voyage
// @Summary      Create voyage
// @Description  Create a new voyage. The `voyage_number` field must be unique.
// @Tags         Voyages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.Voyage true "Voyage payload"
// @Success      201 {object} models.Voyage
// @Failure      400 {object} map[string]string "Bad request"
// @Failure      409 {object} map[string]string "Voyage number already exists"
// @Router       /api/voyages [post]
func CreateVoyage(c *gin.Context) {
	var voyage models.Voyage
	if err := c.ShouldBindJSON(&voyage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate voyage number
	var existing models.Voyage
	if err := config.DB.Where("voyage_number = ?", voyage.VoyageNumber).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Voyage number already exists"})
		return
	}

	if err := config.DB.Create(&voyage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, voyage)
}

// UpdateVoyage handles the PUT request to update an existing voyage
// @Summary      Update voyage
// @Description  Update an existing voyage by ID.
// @Tags         Voyages
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Voyage ID"
// @Param        request body models.Voyage true "Voyage payload"
// @Success      200 {object} map[string]string "Voyage updated successfully"
// @Failure      400 {object} map[string]string "Bad request"
// @Failure      409 {object} map[string]string "Voyage number already exists"
// @Router       /api/voyages/{id} [put]
func UpdateVoyage(c *gin.Context) {
	id := c.Param("id")
	var voyage models.Voyage
	if err := c.ShouldBindJSON(&voyage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check for duplicate voyage number (excluding current ID)
	var existing models.Voyage
	if err := config.DB.Where("voyage_number = ? AND id <> ?", voyage.VoyageNumber, id).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Voyage number already exists"})
		return
	}

	if err := config.DB.Where("id = ?", id).Updates(&voyage).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Voyage updated successfully"})
}

// DeleteVoyage handles the DELETE request to remove a voyage
// @Summary      Delete voyage
// @Description  Delete a voyage by ID.
// @Tags         Voyages
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Voyage ID"
// @Success      200 {object} map[string]string "Voyage deleted successfully"
// @Failure      500 {object} map[string]string "Internal server error"
// @Router       /api/voyages/{id} [delete]
func DeleteVoyage(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Voyage{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Voyage deleted successfully"})
}
