package storage

import (
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateDepresiasi - Menambahkan Depresiasi baru
func CreateDepresiasi(c *gin.Context) {
	var depresiasi types.Depresiasi
	if err := c.ShouldBindJSON(&depresiasi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	result := db.Create(&depresiasi)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, depresiasi)
}

// GetDepresiasiByID - Mendapatkan Depresiasi berdasarkan ID
func GetDepresiasiByID(c *gin.Context) {
	id := c.Param("id")

	var depresiasi types.Depresiasi
	db := database.GetDB()
	result := db.First(&depresiasi, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Depresiasi not found"})
		return
	}

	c.JSON(http.StatusOK, depresiasi)
}

// GetAllDepresiasi - Mendapatkan semua data Depresiasi
func GetAllDepresiasi(c *gin.Context) {
	var depresiasi []types.Depresiasi
	db := database.GetDB()
	
	// Mengambil semua data depresiasi
	if err := db.Find(&depresiasi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Jika tidak ada data, kembalikan array kosong bukan error
	if len(depresiasi) == 0 {
		c.JSON(http.StatusOK, []types.Depresiasi{})
		return
	}
	
	c.JSON(http.StatusOK, depresiasi)
}

// UpdateDepresiasi - Mengupdate data Depresiasi
func UpdateDepresiasi(c *gin.Context) {
	id := c.Param("id")
	var depresiasi types.Depresiasi
	if err := c.ShouldBindJSON(&depresiasi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	result := db.Model(&types.Depresiasi{}).Where("id = ?", id).Updates(depresiasi)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, depresiasi)
}

// DeleteDepresiasi - Menghapus Depresiasi berdasarkan ID
func DeleteDepresiasi(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	result := db.Delete(&types.Depresiasi{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Deleted successfully"})
}
