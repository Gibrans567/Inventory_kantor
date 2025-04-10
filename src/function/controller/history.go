package controller

import (
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Create a new history
func CreateHistory(c *gin.Context) {
	var history types.History
	if err := c.ShouldBindJSON(&history); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set creation time
	history.CreatedAt = time.Now()

	// Create record
	db := database.GetDB()
	result := db.Create(&history)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create history"})
		return
	}

	c.JSON(http.StatusCreated, history)
}

// Get all histories
func GetAllHistories(c *gin.Context) {
	var histories []types.History
	db := database.GetDB()
	result := db.Find(&histories)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch histories"})
		return
	}

	c.JSON(http.StatusOK, histories)
}

// Get history by ID
func GetHistoryByID(c *gin.Context) {
	id := c.Param("id")
	var history types.History

	db := database.GetDB()
	result := db.First(&history, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "History not found"})
		return
	}

	c.JSON(http.StatusOK, history)
}

// Update history
func UpdateHistory(c *gin.Context) {
	id := c.Param("id")
	var history types.History
	db := database.GetDB()

	// Check if history exists
	if err := db.First(&history, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "History not found"})
		return
	}

	// Bind incoming JSON to history struct
	var updateData types.History
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	history.Kategori = updateData.Kategori
	history.Keterangan = updateData.Keterangan
	// Note: We don't update CreatedAt

	// Save changes
	db.Save(&history)

	c.JSON(http.StatusOK, history)
}

// Delete history
func DeleteHistory(c *gin.Context) {
	id := c.Param("id")
	var history types.History
	db := database.GetDB()

	// Check if history exists
	if err := db.First(&history, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "History not found"})
		return
	}

	// Delete the history
	db.Delete(&history)

	c.JSON(http.StatusOK, gin.H{"message": "History deleted successfully"})
}