package controller

import (
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"fmt"
)

// CreateSebaranBarang - Menambahkan SebaranBarang baru
func CreateSebaranBarang(c *gin.Context) {
	var sebaranBarang types.SebaranBarang

	// Binding JSON ke struct
	if err := c.ShouldBindJSON(&sebaranBarang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah IdDivisi valid
	var divisi types.Divisi
	db := database.GetDB()
	if err := db.First(&divisi, sebaranBarang.IdDivisi).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Divisi tidak ditemukan"})
		return
	}

	// Cek apakah IdBarang valid
	var barang types.Inventaris
	if err := db.First(&barang, sebaranBarang.IdBarang).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Barang tidak ditemukan"})
		return
	}

	// Cek apakah IdUser valid
	var user types.User
	if err := db.First(&user, sebaranBarang.IdUser).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User tidak ditemukan"})
		return
	}

	// Cek apakah qty_tersedia di inventaris cukup
	if barang.QtyTersedia < sebaranBarang.QtyBarang {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Qty Barang yang tersedia tidak cukup"})
		return
	}

	// Mulai transaksi database
	tx := db.Begin()

	// Update qty_tersedia dan qty_digunakan pada Inventaris
	barang.QtyTersedia -= sebaranBarang.QtyBarang
	barang.QtyTerpakai += sebaranBarang.QtyBarang
	
	if err := tx.Save(&barang).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate Inventaris"})
		return
	}

	// Simpan data SebaranBarang ke database
	if err := tx.Create(&sebaranBarang).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Buat catatan history untuk Barang Keluar
	now := time.Now()
	
	// Ambil gudang dari barang
	var gudang types.Gudang
	if err := db.First(&gudang, barang.GudangID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data gudang"})
		return
	}
	
	// Membuat pesan history dengan format yang diminta
	historyKeterangan := fmt.Sprintf("Barang %s telah dipindahkan oleh %s dari divisi %s sebanyak %d dari %s ke %s pada %s", 
		barang.NamaBarang,
		user.Name,
		divisi.NamaDivisi,
		sebaranBarang.QtyBarang,
		sebaranBarang.PosisiAwal, // posisi awal (gudang)
		sebaranBarang.PosisiAkhir, // posisi akhir (lokasi baru)
		now.Format("02-01-2006 15:04:05"))
	
	history := types.History{
		Kategori:   "Perpindahan Barang",
		Keterangan: historyKeterangan,
		CreatedAt:  now,
	}
	
	// Simpan data history ke database
	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data history: " + err.Error()})
		return
	}

	// Commit transaksi jika semua operasi berhasil
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data"})
		return
	}

	// Response dengan pesan khusus jika qty_tersedia menjadi 0
	responseData := gin.H{
		"data": sebaranBarang,
		"history": history,
		"message": "Data berhasil disimpan",
	}
	
	if barang.QtyTersedia == 0 {
		responseData["message"] = "Data berhasil disimpan. Semua barang sudah digunakan."
	}

	// Respons sukses dengan data yang telah disimpan
	c.JSON(http.StatusCreated, responseData)
}

// GetSebaranBarangByID - Mendapatkan SebaranBarang berdasarkan ID
func GetSebaranBarangByID(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	var sebaranBarangs []types.SebaranBarang

	// Ambil semua sebaran barang berdasarkan id_barang dan preload relasi
	err := db.
		Preload("Divisi").
		Preload("Inventaris").
		Preload("User").
		Where("id_barang = ?", id).
		Find(&sebaranBarangs).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(sebaranBarangs) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No SebaranBarang found with that id_barang"})
		return
	}

	// Struct untuk response
	type Response struct {
		ID           uint   `json:"id"`
		NamaDivisi   string `json:"nama_divisi"`
		NamaBarang   string `json:"nama_barang"`
		NamaUser     string `json:"nama"`
		QtyBarang    int    `json:"qty_barang"`
		PosisiAwal   string `json:"posisi_awal"`
		PosisiAkhir  string `json:"posisi_akhir"`
		Status       string `json:"status"`
		CreatedAt    time.Time `json:"created_at"`
	}

	// Mapping hasil ke response
	var responseData []Response
	for _, sb := range sebaranBarangs {
		res := Response{
			ID:           sb.ID,
			NamaDivisi:   sb.Divisi.NamaDivisi,
			NamaBarang:   sb.Inventaris.NamaBarang,
			NamaUser:     sb.User.Name,
			QtyBarang:    sb.QtyBarang,
			PosisiAwal:   sb.PosisiAwal,
			PosisiAkhir:  sb.PosisiAkhir,
			Status:       sb.Status,
			CreatedAt:    sb.CreatedAt,
		}
		responseData = append(responseData, res)
	}

	c.JSON(http.StatusOK, responseData)
}

// GetAllSebaranBarang - Mendapatkan semua SebaranBarang dengan hanya menampilkan nama terkait
func GetAllSebaranBarang(c *gin.Context) {
	var sebaranBarangList []types.SebaranBarang
	db := database.GetDB()

	// Mengambil semua data SebaranBarang dengan join ke tabel Divisi, Inventaris, dan User
	result := db.Preload("Divisi").Preload("User").Preload("Inventaris").
		Find(&sebaranBarangList)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan data SebaranBarang"})
		return
	}

	// Membuat response yang hanya berisi nama-nama terkait
	var response []gin.H
	for _, sebaranBarang := range sebaranBarangList {
		response = append(response, gin.H{
			"divisi":     sebaranBarang.Divisi.NamaDivisi, // Asumsi nama divisi ada di struct Divisi
			"barang":     sebaranBarang.Inventaris.NamaBarang, // Asumsi nama barang ada di struct Inventaris
			"user":       sebaranBarang.User.Name, // Asumsi nama user ada di struct User
			"qty_barang": sebaranBarang.QtyBarang,
			"posisi_awal": sebaranBarang.PosisiAwal,
			"posisi_akhir": sebaranBarang.PosisiAkhir,
			"createdAt": sebaranBarang.CreatedAt,
			"updatedAt": sebaranBarang.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response)
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
	
	// Begin transaction
	tx := db.Begin()
	
	// Retrieve the current data before update
	var currentBarang types.SebaranBarang
	if err := tx.Where("id = ?", id).First(&currentBarang).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	
	// Check if there is another record with the same id_barang and same posisi_akhir (excluding current record)
	var targetBarang types.SebaranBarang
	findResult := tx.Where("id_barang = ? AND posisi_akhir = ? AND id != ?", 
		sebaranBarang.IdBarang, sebaranBarang.PosisiAkhir, id).First(&targetBarang)
	
	if findResult.Error == nil {
		// Found a match, add the quantity to the target
		targetBarang.QtyBarang += sebaranBarang.QtyBarang
		
		// Update the target item with the increased quantity
		if err := tx.Model(&targetBarang).Update("qty_barang", targetBarang.QtyBarang).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		// Delete the current item since it's been merged
		if err := tx.Delete(&types.SebaranBarang{}, id).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		
		tx.Commit()
		c.JSON(http.StatusOK, gin.H{
			"message": "Barang successfully merged with existing record",
			"target_item": targetBarang,
		})
		return
	}
	
	// No matching item found, just update as normal
	result := tx.Model(&types.SebaranBarang{}).Where("id = ?", id).Updates(sebaranBarang)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	
	tx.Commit()
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
