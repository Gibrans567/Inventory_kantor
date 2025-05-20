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

	// Ambil data Divisi
	var divisi types.Divisi
	if err := db.First(&divisi, inventaris.DivisiID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Divisi dari barang tidak ditemukan",
			"data":    nil,
		})
		return
	}

	// Ambil data User
	var user types.User
	if err := db.First(&user, pengajuan.IdUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "User tidak ditemukan",
			"data":    nil,
		})
		return
	}

	// Isi field relasi
	pengajuan.StatusKepemilikan = divisi.NamaDivisi
	pengajuan.IdDivisi = inventaris.DivisiID
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

	// Buat response lengkap
	response := gin.H{
		"id":                   pengajuan.ID,
		"id_user":              pengajuan.IdUser,
		"name":                 user.Name,
		"id_barang":            pengajuan.IdBarang,
		"nama_barang":          inventaris.NamaBarang,
		"id_divisi":            pengajuan.IdDivisi,
		"nama_divisi":          divisi.NamaDivisi,
		"tanggal_pengajuan":    pengajuan.TanggalPengajuan,
		"status_permohonan":    pengajuan.StatusPermohonan,
		"status_pengembalian":  pengajuan.StatusPengembalian,
		"status_kepemilikan":   pengajuan.StatusKepemilikan,
		"note":           pengajuan.Note,
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Data berhasil dibuat",
		"data":    response,
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

	for i := range pengajuan {
		pengajuan[i].NamaBarang = pengajuan[i].Inventaris.NamaBarang
		pengajuan[i].NamaUser = pengajuan[i].User.Name
		pengajuan[i].NamaDivisi = pengajuan[i].Divisi.NamaDivisi

		// Ambil data approver user dari db
		var approverUser types.User
		if err := db.First(&approverUser, pengajuan[i].IdApprover).Error; err == nil {
			pengajuan[i].NamaApprover = approverUser.Name
		} else {
			pengajuan[i].NamaApprover = ""
		}
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

	// Mulai transaksi database
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Cek jika status_permohonan berubah menjadi "Disetujui"
	if input.StatusPermohonan == "Disetujui" && existing.StatusPermohonan != "Disetujui" {
		// Ambil data inventaris berdasarkan id (id inventaris = id_barang pengajuan)
		var inventaris types.Inventaris
		if err := tx.Where("id = ?", existing.IdBarang).First(&inventaris).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "Data inventaris tidak ditemukan",
				"data":    nil,
			})
			return
		}

		// Validasi apakah qty_tersedia cukup
		if inventaris.QtyTersedia < existing.QtyBarang {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Stok tersedia tidak mencukupi",
				"data":    nil,
			})
			return
		}

		// Update inventaris: kurangi qty_tersedia dan tambah qty_pinjam
		newQtyTersedia := inventaris.QtyTersedia - existing.QtyBarang
		newQtyPinjam := inventaris.QtyPinjam + existing.QtyBarang

		if err := tx.Model(&inventaris).Updates(map[string]interface{}{
			"qty_tersedia": newQtyTersedia,
			"qty_pinjam":   newQtyPinjam,
		}).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Gagal update inventaris",
				"data":    nil,
			})
			return
		}
	}

	// Update field lainnya
	if input.IdApprover != nil {
		existing.IdApprover = input.IdApprover
	}
	if input.StatusPermohonan != "" {
		existing.StatusPermohonan = input.StatusPermohonan
	}
	if input.QtyBarang != 0 {
		existing.QtyBarang = input.QtyBarang
	}
	if input.Note != "" {
		existing.Note = input.Note
	}
	if input.StatusPengembalian != "" {
		existing.StatusPengembalian = input.StatusPengembalian
	}

	// Simpan perubahan pengajuan
	if err := tx.Save(&existing).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal update data pengajuan",
			"data":    nil,
		})
		return
	}

	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan perubahan",
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

