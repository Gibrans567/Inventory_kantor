package controller

import (
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
	
)

func GetAllMobil(c *gin.Context) {
	db := database.GetDB()
	var mobils []types.Mobil
	
	if err := db.Find(&mobils).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data mobil",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   mobils,
	})
}

// GetMobilByID - GET /mobil/:id
func GetMobilByID(c *gin.Context) {
	db := database.GetDB()
	var mobil types.Mobil
	id := c.Param("id")
	
	if err := db.First(&mobil, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Mobil tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data mobil",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   mobil,
	})
}

// CreateMobil - POST /mobil
func CreateMobil(c *gin.Context) {
	db := database.GetDB()
	var mobil types.Mobil
	
	if err := c.ShouldBindJSON(&mobil); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Format data tidak valid: " + err.Error(),
		})
		return
	}
	
	// Validasi field yang wajib diisi
	if mobil.NamaMobil == "" || mobil.PlatNomor == "" || mobil.TipeMobil == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nama mobil, plat nomor, dan tipe mobil harus diisi",
		})
		return
	}
	
	if err := db.Create(&mobil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menyimpan data mobil",
		})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Mobil berhasil ditambahkan",
		"data":    mobil,
	})
}

// UpdateMobil - PUT /mobil/:id
func UpdateMobil(c *gin.Context) {
	db := database.GetDB()
	var mobil types.Mobil
	id := c.Param("id")
	
	// Cek apakah mobil ada
	if err := db.First(&mobil, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Mobil tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data mobil",
		})
		return
	}
	
	var updateData types.Mobil
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Format data tidak valid: " + err.Error(),
		})
		return
	}
	
	// Validasi field yang wajib diisi
	if updateData.NamaMobil == "" || updateData.PlatNomor == "" || updateData.TipeMobil == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nama mobil, plat nomor, dan tipe mobil harus diisi",
		})
		return
	}
	
	// Update data
	mobil.NamaMobil = updateData.NamaMobil
	mobil.PlatNomor = updateData.PlatNomor
	mobil.TipeMobil = updateData.TipeMobil
	
	if err := db.Save(&mobil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengupdate data mobil",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Mobil berhasil diupdate",
		"data":    mobil,
	})
}

// DeleteMobil - DELETE /mobil/:id
func DeleteMobil(c *gin.Context) {
	db := database.GetDB()
	var mobil types.Mobil
	id := c.Param("id")
	
	// Cek apakah mobil ada
	if err := db.First(&mobil, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Mobil tidak ditemukan",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal mengambil data mobil",
		})
		return
	}
	
	if err := db.Delete(&mobil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Gagal menghapus data mobil",
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Mobil berhasil dihapus",
	})
}
