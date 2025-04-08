package controller

import (
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
	"time"
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
		Perbulan:        hargaDepresiasi,  // Set default ke 1
		Tahun:           hargaDepresiasi*12,  // Set default ke 1
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

// GetInventarisById - Mendapatkan data Inventaris berdasarkan ID dan data SebaranBarang terkait
func GetInventarisById(c *gin.Context) {
    // Mendapatkan ID dari parameter URL
    id := c.Param("id")
    
    // Validasi ID
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak boleh kosong"})
        return
    }
    
    var inventaris types.Inventaris
    db := database.GetDB()
    
    // Mengambil data inventaris berdasarkan ID dan menampilkan nama Gudang dan Kategori
    if err := db.Preload("Gudang", func(db *gorm.DB) *gorm.DB {
        return db.Select("id, nama_gudang")  // Hanya mengambil id dan nama_gudang dari Gudang
    }).Preload("Kategori", func(db *gorm.DB) *gorm.DB {
        return db.Select("id, nama_kategori")  // Hanya mengambil id dan nama_kategori dari Kategori
    }).First(&inventaris, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Data inventaris tidak ditemukan"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    
    // Mengambil semua data sebaran barang yang memiliki id_barang yang sama
    var sebaranBarang []types.SebaranBarang
    if err := db.Where("id_barang = ?", inventaris.ID).Find(&sebaranBarang).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // Format waktu yang konsisten
    createdAt := inventaris.CreatedAt.Format(time.RFC3339)
    updatedAt := inventaris.UpdatedAt.Format(time.RFC3339)
    tanggalPembelian := inventaris.TanggalPembelian.Format(time.RFC3339)
    
    // Menyiapkan hasil yang berisi data inventaris dengan nama gudang dan kategori
    result := gin.H{
        "id":                 inventaris.ID,
        "tanggal_pembelian":  tanggalPembelian,
        "gudang_id":          inventaris.GudangID,
        "gudang_nama":        inventaris.Gudang.NamaGudang,     // Nama Gudang
        "kategori_id":        inventaris.KategoriID,
        "kategori_nama":      inventaris.Kategori.NamaKategori, // Nama Kategori
        "nama_barang":        inventaris.NamaBarang,
        "qty_barang":         inventaris.QtyBarang,
        "qty_terpakai":       inventaris.QtyTerpakai,
        "qty_tersedia":       inventaris.QtyTersedia,
        "harga_pembelian":    inventaris.HargaPembelian,
        "spesifikasi":        inventaris.Spesifikasi,
        "total_nilai":        inventaris.TotalNilai,
        "upload_nota":        inventaris.UploadNota,
        "created_at":         createdAt,
        "updated_at":         updatedAt,
        "sebaran_barang":     sebaranBarang,                    // Data sebaran barang di akhir
    }

    c.JSON(http.StatusOK, result)
}

func GetInventarisByDate(c *gin.Context) {
	// Struct untuk menerima body request
	var request struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}

	// Binding JSON
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format JSON tidak valid"})
		return
	}

	// Parse tanggal
	startDate, err := time.Parse("2006-01-02", request.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date harus format YYYY-MM-DD"})
		return
	}
	endDate, err := time.Parse("2006-01-02", request.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date harus format YYYY-MM-DD"})
		return
	}
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second) // Biar sampai akhir hari

	// Query database
	db := database.GetDB()
	query := db.Preload("Gudang", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, nama_gudang")
	}).Preload("Kategori", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, nama_kategori")
	}).Where("created_at BETWEEN ? AND ?", startDate, endDate)

	var inventarisList []types.Inventaris
	if err := query.Find(&inventarisList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []gin.H
	for _, inv := range inventarisList {
		var sebaranBarang []types.SebaranBarang
		if err := db.Where("id_barang = ? AND created_at BETWEEN ? AND ?", inv.ID, startDate, endDate).Find(&sebaranBarang).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		result = append(result, gin.H{
			"id":                 inv.ID,
			"tanggal_pembelian": inv.TanggalPembelian.Format(time.RFC3339),
			"gudang_id":          inv.GudangID,
			"gudang_nama":        inv.Gudang.NamaGudang,
			"kategori_id":        inv.KategoriID,
			"kategori_nama":      inv.Kategori.NamaKategori,
			"nama_barang":        inv.NamaBarang,
			"qty_barang":         inv.QtyBarang,
			"qty_terpakai":       inv.QtyTerpakai,
			"qty_tersedia":       inv.QtyTersedia,
			"harga_pembelian":    inv.HargaPembelian,
			"spesifikasi":        inv.Spesifikasi,
			"total_nilai":        inv.TotalNilai,
			"upload_nota":        inv.UploadNota,
			"created_at":         inv.CreatedAt.Format(time.RFC3339),
			"updated_at":         inv.UpdatedAt.Format(time.RFC3339),
			"sebaran_barang":     sebaranBarang,
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
