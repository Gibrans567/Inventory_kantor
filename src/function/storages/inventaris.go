package storage

import (
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
)

// CreateInventaris - Menambahkan Inventaris baru dan mencatat depresiasi
func CreateInventaris(c *gin.Context) {
	var inventaris types.Inventaris
	db := database.GetDB()
	// Bind JSON request ke struct Inventaris
	if err := c.ShouldBindJSON(&inventaris); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Mengecek apakah GudangID ada dalam tabel Gudang
	var gudangExists bool
	
	if err := db.Model(&types.Gudang{}).Where("id = ?", inventaris.GudangID).First(&types.Gudang{}).Error; err != nil {
		gudangExists = false
	} else {
		gudangExists = true
	}

	// Mengecek apakah KategoriID ada dalam tabel Kategori
	var kategoriExists bool
	if err := db.Model(&types.Kategori{}).Where("id = ?", inventaris.KategoriID).First(&types.Kategori{}).Error; err != nil {
		kategoriExists = false
	} else {
		kategoriExists = true
	}

	// Jika GudangID atau KategoriID tidak ada
	if !gudangExists || !kategoriExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "GudangID atau KategoriID tidak valid. Pastikan data ada di tabel Gudang atau Kategori.",
		})
		return
	}

	// Mengecualikan input untuk qty_terpakai
	inventaris.QtyTerpakai = 0 // Mengatur qty_terpakai ke 0 jika tidak diinput

	// Menghitung total_nilai berdasarkan harga_pembelian dan qty_barang
	inventaris.TotalNilai = inventaris.HargaPembelian * float64(inventaris.QtyBarang)

	// Mengatur qty_tersedia sama dengan qty_barang
	inventaris.QtyTersedia = inventaris.QtyBarang

	// Memulai transaksi database
	tx := db.Begin()
	
	// Menyimpan data inventaris ke database
	if err := tx.Create(&inventaris).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Menghitung depresiasi (2.5% dari harga pembelian)
	hargaDepresiasi := int(inventaris.HargaPembelian * 0.025)
	
	// Membuat record depresiasi
	depresiasi := types.Depresiasi{
		IdGudang:        inventaris.GudangID,
		IdBarang:        inventaris.ID,
		HargaDepresiasi: hargaDepresiasi,
		Perbulan:        1,  // Set default ke 1
		Tahun:           1,  // Set default ke 1
	}
	
	// Menyimpan data depresiasi ke database
	if err := tx.Create(&depresiasi).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data depresiasi: " + err.Error()})
		return
	}
	
	// Commit transaksi
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal melakukan commit transaksi: " + err.Error()})
		return
	}

	// Mengirimkan respon jika berhasil
	c.JSON(http.StatusCreated, gin.H{
		"inventaris": inventaris,
		"depresiasi": depresiasi,
	})
}


// GetAllInventaris - Mendapatkan semua data Inventaris
func GetAllInventaris(c *gin.Context) {
    var inventaris []types.Inventaris
    db := database.GetDB()
    
    // Mengambil data inventaris dan menampilkan nama Gudang dan Kategori
    if err := db.Preload("Gudang", func(db *gorm.DB) *gorm.DB {
        return db.Select("id, nama_gudang")  // Hanya mengambil id dan nama_gudang dari Gudang
    }).Preload("Kategori", func(db *gorm.DB) *gorm.DB {
        return db.Select("id, nama_kategori")  // Hanya mengambil id dan nama_kategori dari Kategori
    }).Find(&inventaris).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Menyiapkan hasil yang berisi semua data inventaris dengan nama gudang dan kategori
    var result []gin.H
    for _, item := range inventaris {
        result = append(result, gin.H{
            "id":                 item.ID,
            "tanggal_pembelian":  item.TanggalPembelian,
            "gudang_id":          item.GudangID,
            "kategori_id":        item.KategoriID,
            "nama_barang":        item.NamaBarang,
            "qty_barang":         item.QtyBarang,
            "harga_pembelian":    item.HargaPembelian,
            "spesifikasi":        item.Spesifikasi,
            "qty_tersedia":       item.QtyTersedia,
            "qty_terpakai":       item.QtyTerpakai,
            "total_nilai":        item.TotalNilai,
            "upload_nota":        item.UploadNota,
            "created_at":         item.CreatedAt,
            "updated_at":         item.UpdatedAt,
            "gudang_nama":        item.Gudang.NamaGudang,     // Nama Gudang
            "kategori_nama":      item.Kategori.NamaKategori, // Nama Kategori
        })
    }

    c.JSON(http.StatusOK, result)
}



// UpdateInventaris - Mengupdate data Inventaris
func UpdateInventaris(c *gin.Context) {
	id := c.Param("id")
	var inventaris types.Inventaris
	if err := c.ShouldBindJSON(&inventaris); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	result := db.Model(&types.Inventaris{}).Where("id = ?", id).Updates(inventaris)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, inventaris)
}

// DeleteInventaris - Menghapus Inventaris berdasarkan ID
func DeleteInventaris(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	result := db.Delete(&types.Inventaris{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Deleted successfully"})
}
