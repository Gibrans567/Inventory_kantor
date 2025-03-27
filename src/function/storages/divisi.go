package storage

import (
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateDivisi - Menambahkan Divisi baru
func CreateDivisi(c *gin.Context) {
	var divisi types.Divisi
	if err := c.ShouldBindJSON(&divisi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	result := db.Create(&divisi)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, divisi)
}

// GetDivisiByID - Mendapatkan Divisi berdasarkan ID
func GetDivisiByID(c *gin.Context) {
	id := c.Param("id")

	var divisi types.Divisi
	db := database.GetDB()
	result := db.First(&divisi, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Divisi not found"})
		return
	}

	c.JSON(http.StatusOK, divisi)
}

// UpdateDivisi - Mengupdate data Divisi
func UpdateDivisi(c *gin.Context) {
	id := c.Param("id")
	var divisi types.Divisi
	if err := c.ShouldBindJSON(&divisi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	result := db.Model(&types.Divisi{}).Where("id = ?", id).Updates(divisi)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, divisi)
}

// DeleteDivisi - Menghapus Divisi berdasarkan ID
func DeleteDivisi(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	result := db.Delete(&types.Divisi{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Deleted successfully"})
}
