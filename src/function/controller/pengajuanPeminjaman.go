package controller

import (
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	
)

func CreatePengajuan(c *gin.Context) {
	db := database.GetDB()

	var pengajuan types.PengajuanPeminjaman
	if err := c.ShouldBindJSON(&pengajuan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// Ambil data Inventaris berdasarkan id_barang
	var inventaris types.Inventaris
	if err := db.First(&inventaris, pengajuan.IdBarang).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Data barang tidak ditemukan",
			"data":    nil,
		})
		return
	}

	// Ambil nama_divisi berdasarkan id_divisi dari Inventaris
	var divisi types.Divisi
	if err := db.First(&divisi, inventaris.DivisiID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Divisi dari barang tidak ditemukan",
			"data":    nil,
		})
		return
	}

	// Isi field yang tergantung data relasi
	pengajuan.StatusKepemilikan = divisi.NamaDivisi
	pengajuan.IdDivisi = inventaris.DivisiID // opsional: agar tetap konsisten
	pengajuan.TanggalPengajuan = time.Now()

	if pengajuan.StatusPermohonan == "" {
		pengajuan.StatusPermohonan = "Menunggu Approve"
	}

	if pengajuan.StatusPengembalian == "" {
		pengajuan.StatusPengembalian = "Belum Dikembalikan"
	}

	// Simpan ke database
	if err := db.Create(&pengajuan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal membuat data",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Data berhasil dibuat",
		"data":    pengajuan,
	})
}

func GetAllPengajuan(c *gin.Context) {
	db := database.GetDB()

	var pengajuan []types.PengajuanPeminjaman
	if err := db.
		Preload("Inventaris").
		Preload("User").
		Preload("Divisi").
		Find(&pengajuan).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch data",
			"data":    nil,
		})
		return
	}

	// Mengisi nama_barang, nama_user, nama_divisi secara manual dari relasi
	for i := range pengajuan {
		pengajuan[i].NamaBarang = pengajuan[i].Inventaris.NamaBarang
		pengajuan[i].NamaUser = pengajuan[i].User.Name
		pengajuan[i].NamaDivisi = pengajuan[i].Divisi.NamaDivisi
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data fetched successfully",
		"data":    pengajuan,
	})
}

func GetPengajuanByID(c *gin.Context) {
	db := database.GetDB()

	id := c.Param("id")
	var pengajuan types.PengajuanPeminjaman
	if err := db.First(&pengajuan, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data not found",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data found",
		"data":    pengajuan,
	})
}


func UpdatePengajuan(c *gin.Context) {
	db := database.GetDB()
	id := c.Param("id")

	// Ambil data lama dari database
	var existing types.PengajuanPeminjaman
	if err := db.First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Data tidak ditemukan",
			"data":    nil,
		})
		return
	}

	// Binding data baru dari request body
	var input types.PengajuanPeminjaman
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// Langsung update id_approver dan status_permohonan
	if input.IdApprover != nil {
		existing.IdApprover = input.IdApprover
	}
	if input.StatusPermohonan != "" {
		existing.StatusPermohonan = input.StatusPermohonan
	}

	// Update field lainnya jika ada perubahan
	if input.QtyBarang != 0 {
		existing.QtyBarang = input.QtyBarang
	}
	if input.Note != "" {
		existing.Note = input.Note
	}
	if input.PosisiAkhir != "" {
		existing.PosisiAkhir = input.PosisiAkhir
	}
	if input.StatusPengembalian != "" {
		existing.StatusPengembalian = input.StatusPengembalian
	}

	// Simpan perubahan
	if err := db.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal update data",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data berhasil diperbarui",
		"data":    existing,
	})
}


func DeletePengajuan(c *gin.Context) {
	db := database.GetDB()

	id := c.Param("id")
	if err := db.Delete(&types.PengajuanPeminjaman{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete data",
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Data deleted successfully",
		"data":    nil,
	})
}

