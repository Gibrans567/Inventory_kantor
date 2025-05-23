package controller

import (
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
	
)

// GetAllPeminjamanMobil - GET /peminjaman-mobil
func GetAllPeminjamanMobil(c *gin.Context) {
	db := database.GetDB()
	var peminjamans []types.PeminjamanMobil
	
	if err := db.Preload("Mobil").Preload("User").Preload("Divisi").Find(&peminjamans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data peminjaman mobil",
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data peminjaman mobil berhasil diambil",
		"data":    peminjamans,
	})
}

// GetPeminjamanMobilByID - GET /peminjaman-mobil/:id
func GetPeminjamanMobilByID(c *gin.Context) {
	db := database.GetDB()
	var peminjaman types.PeminjamanMobil
	id := c.Param("id")
	
	if err := db.Preload("Mobil").Preload("User").Preload("Divisi").First(&peminjaman, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Peminjaman mobil tidak ditemukan",
				"data":    nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data peminjaman mobil",
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data peminjaman mobil berhasil diambil",
		"data":    peminjaman,
	})
}

// CreatePeminjamanMobil - POST /peminjaman-mobil
func CreatePeminjamanMobil(c *gin.Context) {
	db := database.GetDB()
	var peminjaman types.PeminjamanMobil
	
	if err := c.ShouldBindJSON(&peminjaman); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format data tidak valid: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	// Validasi field yang wajib diisi
	if peminjaman.IdMobil == 0 || peminjaman.IdUser == 0 || peminjaman.IdDivisi == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID Mobil, ID User, dan ID Divisi harus diisi",
			"data":    nil,
		})
		return
	}
	
	if peminjaman.TanggalPinjam.IsZero() || peminjaman.TanggalKembali.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Tanggal pinjam dan tanggal kembali harus diisi",
			"data":    nil,
		})
		return
	}
	
	
	
	// Validasi tanggal kembali tidak boleh sebelum tanggal pinjam
	if peminjaman.TanggalKembali.Before(peminjaman.TanggalPinjam) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Tanggal kembali tidak boleh sebelum tanggal pinjam",
			"data":    nil,
		})
		return
	}
	
	// Cek apakah mobil, user, dan divisi exist
	var mobil types.Mobil
	var user types.User
	var divisi types.Divisi
	
	if err := db.First(&mobil, peminjaman.IdMobil).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Mobil tidak ditemukan",
			"data":    nil,
		})
		return
	}
	
	if err := db.First(&user, peminjaman.IdUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User tidak ditemukan",
			"data":    nil,
		})
		return
	}
	
	if err := db.First(&divisi, peminjaman.IdDivisi).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Divisi tidak ditemukan",
			"data":    nil,
		})
		return
	}
	
	// Cek apakah mobil sedang dipinjam
	var existingPeminjaman types.PeminjamanMobil
	result := db.Where("id_mobil = ? AND status_pinjam IN (?, ?)", 
		peminjaman.IdMobil, "dipinjam", "pending").First(&existingPeminjaman)
	
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Mobil sedang dipinjam atau dalam status pending",
			"data":    nil,
		})
		return
	}
	
	if err := db.Create(&peminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan data peminjaman mobil",
			"data":    nil,
		})
		return
	}
	
	// Load relasi untuk response
	db.Preload("Mobil").Preload("User").Preload("Divisi").First(&peminjaman, peminjaman.ID)
	
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Peminjaman mobil berhasil ditambahkan",
		"data":    peminjaman,
	})
}

// UpdatePeminjamanMobil - PUT /peminjaman-mobil/:id
func UpdatePeminjamanMobil(c *gin.Context) {
	db := database.GetDB()
	var peminjaman types.PeminjamanMobil
	id := c.Param("id")
	
	// Cek apakah peminjaman ada
	if err := db.First(&peminjaman, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Peminjaman mobil tidak ditemukan",
				"data":    nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data peminjaman mobil",
			"data":    nil,
		})
		return
	}
	
	var updateData types.PeminjamanMobil
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Format data tidak valid: " + err.Error(),
			"data":    nil,
		})
		return
	}
	
	// Validasi field yang wajib diisi
	if updateData.IdMobil == 0 || updateData.IdUser == 0 || updateData.IdDivisi == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "ID Mobil, ID User, dan ID Divisi harus diisi",
			"data":    nil,
		})
		return
	}
	
	if updateData.TanggalPinjam.IsZero() || updateData.TanggalKembali.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Tanggal pinjam dan tanggal kembali harus diisi",
			"data":    nil,
		})
		return
	}
	
	
	
	// Validasi tanggal kembali tidak boleh sebelum tanggal pinjam
	if updateData.TanggalKembali.Before(updateData.TanggalPinjam) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Tanggal kembali tidak boleh sebelum tanggal pinjam",
			"data":    nil,
		})
		return
	}
	
	// Cek apakah mobil, user, dan divisi exist (jika diubah)
	if updateData.IdMobil != peminjaman.IdMobil {
		var mobil types.Mobil
		if err := db.First(&mobil, updateData.IdMobil).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Mobil tidak ditemukan",
				"data":    nil,
			})
			return
		}
		
		// Cek apakah mobil baru sedang dipinjam
		var existingPeminjaman types.PeminjamanMobil
		result := db.Where("id_mobil = ? AND status_pinjam IN (?, ?) AND id != ?", 
			updateData.IdMobil, "dipinjam", "pending", peminjaman.ID).First(&existingPeminjaman)
		
		if result.Error == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Mobil sedang dipinjam atau dalam status pending",
				"data":    nil,
			})
			return
		}
	}
	
	if updateData.IdUser != peminjaman.IdUser {
		var user types.User
		if err := db.First(&user, updateData.IdUser).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "User tidak ditemukan",
				"data":    nil,
			})
			return
		}
	}
	
	if updateData.IdDivisi != peminjaman.IdDivisi {
		var divisi types.Divisi
		if err := db.First(&divisi, updateData.IdDivisi).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Divisi tidak ditemukan",
				"data":    nil,
			})
			return
		}
	}
	
	// Update data
	peminjaman.IdMobil = updateData.IdMobil
	peminjaman.IdUser = updateData.IdUser
	peminjaman.IdDivisi = updateData.IdDivisi
	peminjaman.TanggalPinjam = updateData.TanggalPinjam
	peminjaman.TanggalKembali = updateData.TanggalKembali
	
	if err := db.Save(&peminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengupdate data peminjaman mobil",
			"data":    nil,
		})
		return
	}
	
	// Load relasi untuk response
	db.Preload("Mobil").Preload("User").Preload("Divisi").First(&peminjaman, peminjaman.ID)
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Peminjaman mobil berhasil diupdate",
		"data":    peminjaman,
	})
}

// DeletePeminjamanMobil - DELETE /peminjaman-mobil/:id
func DeletePeminjamanMobil(c *gin.Context) {
	db := database.GetDB()
	var peminjaman types.PeminjamanMobil
	id := c.Param("id")
	
	// Cek apakah peminjaman ada
	if err := db.First(&peminjaman, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Peminjaman mobil tidak ditemukan",
				"data":    nil,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data peminjaman mobil",
			"data":    nil,
		})
		return
	}
	
	if err := db.Delete(&peminjaman).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menghapus data peminjaman mobil",
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Peminjaman mobil berhasil dihapus",
		"data":    nil,
	})
}

// GetPeminjamanByStatus - GET /peminjaman-mobil/status/:status
func GetPeminjamanByStatus(c *gin.Context) {
	db := database.GetDB()
	var peminjamans []types.PeminjamanMobil
	status := c.Param("status")
	
	if err := db.Preload("Mobil").Preload("User").Preload("Divisi").
		Where("status_pinjam = ?", status).Find(&peminjamans).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal mengambil data peminjaman mobil",
			"data":    nil,
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data peminjaman mobil berhasil diambil berdasarkan status",
		"data":    peminjamans,
	})
}