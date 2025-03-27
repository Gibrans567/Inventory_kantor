package storage

import (
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateSebaranBarang - Menambahkan SebaranBarang baru
func CreateSebaranBarang(c *gin.Context) {
	var sebaranBarang types.SebaranBarang
	if err := c.ShouldBindJSON(&sebaranBarang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	result := db.Create(&sebaranBarang)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, sebaranBarang)
}

// GetSebaranBarangByID - Mendapatkan SebaranBarang berdasarkan ID
func GetSebaranBarangByID(c *gin.Context) {
	id := c.Param("id")

	var sebaranBarang types.SebaranBarang
	db := database.GetDB()
	result := db.First(&sebaranBarang, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "SebaranBarang not found"})
		return
	}

	c.JSON(http.StatusOK, sebaranBarang)
}

// UpdateSebaranBarang - Mengupdate data SebaranBarang
func UpdateSebaranBarang(c *gin.Context) {
	id := c.Param("id")
	var sebaranBarang types.SebaranBarang
	if err := c.ShouldBindJSON(&sebaranBarang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	result := db.Model(&types.SebaranBarang{}).Where("id = ?", id).Updates(sebaranBarang)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, sebaranBarang)
}

// DeleteSebaranBarang - Menghapus SebaranBarang berdasarkan ID
func DeleteSebaranBarang(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	result := db.Delete(&types.SebaranBarang{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Deleted successfully"})
}
