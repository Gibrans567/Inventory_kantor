package storage

import (
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateGudang - Menambahkan Gudang baru
func CreateGudang(c *gin.Context) {
	var gudang types.Gudang
	if err := c.ShouldBindJSON(&gudang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	result := db.Create(&gudang)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gudang)
}

// GetGudangByID - Mendapatkan Gudang berdasarkan ID
func GetGudangByID(c *gin.Context) {
	id := c.Param("id")

	var gudang types.Gudang
	db := database.GetDB()
	result := db.First(&gudang, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gudang not found"})
		return
	}

	c.JSON(http.StatusOK, gudang)
}

// UpdateGudang - Mengupdate data Gudang
func UpdateGudang(c *gin.Context) {
	id := c.Param("id")
	var gudang types.Gudang
	if err := c.ShouldBindJSON(&gudang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	result := db.Model(&types.Gudang{}).Where("id = ?", id).Updates(gudang)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gudang)
}

// DeleteGudang - Menghapus Gudang berdasarkan ID
func DeleteGudang(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	result := db.Delete(&types.Gudang{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Deleted successfully"})
}
