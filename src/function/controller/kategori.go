package controller

import (
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateKategori - Menambahkan Kategori baru
func CreateKategori(c *gin.Context) {
	var kategori types.Kategori
	if err := c.ShouldBindJSON(&kategori); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	result := db.Create(&kategori)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, kategori)
}

// GetKategoriByID - Mendapatkan Kategori berdasarkan ID
func GetKategoriByID(c *gin.Context) {
	id := c.Param("id")

	var kategori types.Kategori
	db := database.GetDB()
	result := db.First(&kategori, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kategori not found"})
		return
	}

	c.JSON(http.StatusOK, kategori)
}

// UpdateKategori - Mengupdate data Kategori
func GetAllKategori(c *gin.Context) {
	db := database.GetDB()
	var kategoris []types.Kategori

	result := db.Find(&kategoris)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": kategoris})
}

func UpdateKategori(c *gin.Context) {
	id := c.Param("id")
	var kategori types.Kategori
	if err := c.ShouldBindJSON(&kategori); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	result := db.Model(&types.Kategori{}).Where("id = ?", id).Updates(kategori)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, kategori)
}

// DeleteKategori - Menghapus Kategori berdasarkan ID
func DeleteKategori(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	result := db.Delete(&types.Kategori{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Deleted successfully"})
}
