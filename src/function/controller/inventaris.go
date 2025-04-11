package controller

import (
	"inventory/src/types"
	"inventory/src/function/database"

	"github.com/gin-gonic/gin"
	"net/http"
	"gorm.io/gorm"
	"time"
	"fmt"

)

// CreateInventaris - Menambahkan Inventaris baru dan mencatat depresiasi
func CreateInventaris(c *gin.Context) {
    var inventaris types.Inventaris
    db := database.GetDB()

    if err := c.ShouldBindJSON(&inventaris); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Validation code remains the same...
    var gudang types.Gudang
    if err := db.Where("id = ?", inventaris.GudangID).First(&gudang).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "GudangID tidak valid"})
        return
    }

    var kategori types.Kategori
    if err := db.Where("id = ?", inventaris.KategoriID).First(&kategori).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "KategoriID tidak valid"})
        return
    }

	var divisi types.Divisi
    if err := db.Where("id = ?", inventaris.DivisiID).First(&divisi).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "KategoriID tidak valid"})
        return
    }

	var user types.User
    if err := db.Where("id = ?", inventaris.UserID).First(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "KategoriID tidak valid"})
        return
    }

    inventaris.QtyTerpakai = 0
    inventaris.TotalNilai = inventaris.HargaPembelian * (inventaris.QtyBarang)
    inventaris.QtyTersedia = inventaris.QtyBarang

    tx := db.Begin()
    if err := tx.Create(&inventaris).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	hargaDepresiasi := int(float64(inventaris.HargaPembelian) * 0.025)
    depresiasi := types.Depresiasi{
        IdGudang:        inventaris.GudangID,
        IdBarang:        inventaris.ID,
        HargaDepresiasi: hargaDepresiasi,
        Perbulan:        hargaDepresiasi,
        Tahun:           hargaDepresiasi * 12,
    }

    if err := tx.Create(&depresiasi).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan depresiasi: " + err.Error()})
        return
    }

    now := time.Now()
    history := types.History{
        Kategori:   "Barang Masuk",
        Keterangan: fmt.Sprintf("Pada %s barang %s telah masuk ke gudang %s oleh %s dari Divisi %s ", now.Format("02-01-2006 15:04:05"), inventaris.NamaBarang, gudang.NamaGudang, user.Name,divisi.NamaDivisi),
        CreatedAt:  now,
    }

    if err := tx.Create(&history).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan history: " + err.Error()})
        return
    }

    // Membuat Jadwal Depresiasi
    nextRun := time.Now().AddDate(0, 1, 0)
    jadwal := types.JadwalDepresiasi{
        IdBarang: inventaris.ID,
        NextRun:  nextRun,
    }
    if err := tx.Create(&jadwal).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal simpan jadwal depresiasi: " + err.Error()})
        return
    }

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal commit: " + err.Error()})
        return
    }



    c.JSON(http.StatusCreated, gin.H{
        "inventaris": inventaris,
        "depresiasi": depresiasi,
        "history":    history,
        "jadwal":     jadwal,
    })
}

// GetAllInventaris - Mendapatkan semua data Inventaris
func GetAllInventaris(c *gin.Context) {
    var inventaris []types.Inventaris
    db := database.GetDB()
    
    // Mengambil data inventaris dan menampilkan nama Gudang, Kategori, dan Divisi
    if err := db.Preload("Gudang", func(db *gorm.DB) *gorm.DB {
        return db.Select("id, nama_gudang")  // Hanya mengambil id dan nama_gudang dari Gudang
    }).Preload("Kategori", func(db *gorm.DB) *gorm.DB {
        return db.Select("id, nama_kategori")  // Hanya mengambil id dan nama_kategori dari Kategori
    }).Preload("Divisi", func(db *gorm.DB) *gorm.DB {
        return db.Select("id, nama_divisi")  // Hanya mengambil id dan nama_divisi dari Divisi
    }).Find(&inventaris).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Menyiapkan hasil yang berisi semua data inventaris dengan nama gudang, kategori, dan divisi
    var result []gin.H
    for _, item := range inventaris {
        result = append(result, gin.H{
            "id":                 item.ID,
            "tanggal_pembelian":  item.TanggalPembelian,
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
            "divisi_nama":        item.Divisi.NamaDivisi,     // Nama Divisi
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
    
    // Mengambil data inventaris berdasarkan ID dan menampilkan nama Gudang, Kategori, dan Divisi
    if err := db.Preload("Gudang", func(db *gorm.DB) *gorm.DB {
        return db.Select("id, nama_gudang")  // Hanya mengambil id dan nama_gudang dari Gudang
    }).Preload("Kategori", func(db *gorm.DB) *gorm.DB {
        return db.Select("id, nama_kategori")  // Hanya mengambil id dan nama_kategori dari Kategori
    }).Preload("Divisi", func(db *gorm.DB) *gorm.DB {
        return db.Select("id, nama_divisi")  // Hanya mengambil id dan nama_divisi dari Divisi
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
    
    // Menyiapkan hasil yang berisi data inventaris dengan nama gudang, kategori, dan divisi
    result := gin.H{
        "id":                 inventaris.ID,
        "tanggal_pembelian":  tanggalPembelian,
        "gudang_id":          inventaris.GudangID,
        "gudang_nama":        inventaris.Gudang.NamaGudang,     // Nama Gudang
        "kategori_id":        inventaris.KategoriID,
        "kategori_nama":      inventaris.Kategori.NamaKategori, // Nama Kategori
        "divisi_id":          inventaris.DivisiID,              // ID Divisi
        "divisi_nama":        inventaris.Divisi.NamaDivisi,     // Nama Divisi
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
	}).Preload("Divisi", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, nama_divisi")  // Menambahkan preload untuk Divisi
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
			"divisi_id":          inv.DivisiID,              // Menambahkan divisi_id
			"divisi_nama":        inv.Divisi.NamaDivisi,     // Menambahkan divisi_nama
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

func GetInventarisByCategory(c *gin.Context) {
	// Query database
	db := database.GetDB()
	query := db.Preload("Gudang", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, nama_gudang")
	}).Preload("Kategori", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, nama_kategori")
	}).Preload("Divisi", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, nama_divisi")
	})

	var inventarisList []types.Inventaris
	if err := query.Find(&inventarisList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Mengelompokkan hasil berdasarkan kategori
	categoryMap := make(map[uint]gin.H)
	
	for _, inv := range inventarisList {
		var sebaranBarang []types.SebaranBarang
		if err := db.Where("id_barang = ?", inv.ID).Find(&sebaranBarang).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		invData := gin.H{
			"id":                inv.ID,
			"tanggal_pembelian": inv.TanggalPembelian.Format(time.RFC3339),
			"gudang_id":         inv.GudangID,
			"gudang_nama":       inv.Gudang.NamaGudang,
			"divisi_id":         inv.DivisiID,        // Menambahkan divisi_id
			"divisi_nama":       inv.Divisi.NamaDivisi, // Menambahkan divisi_nama
			"nama_barang":       inv.NamaBarang,
			"qty_barang":        inv.QtyBarang,
			"qty_terpakai":      inv.QtyTerpakai,
			"qty_tersedia":      inv.QtyTersedia,
			"harga_pembelian":   inv.HargaPembelian,
			"spesifikasi":       inv.Spesifikasi,
			"total_nilai":       inv.TotalNilai,
			"upload_nota":       inv.UploadNota,
			"created_at":        inv.CreatedAt.Format(time.RFC3339),
			"updated_at":        inv.UpdatedAt.Format(time.RFC3339),
			"sebaran_barang":    sebaranBarang,
		}

		// Jika kategori belum ada di map, buat kategori baru
		if _, exists := categoryMap[inv.KategoriID]; !exists {
			categoryMap[inv.KategoriID] = gin.H{
				"kategori_id":   inv.KategoriID,
				"kategori_nama": inv.Kategori.NamaKategori,
				"items":         []gin.H{invData},
			}
		} else {
			// Jika kategori sudah ada, tambahkan item ke daftar items di kategori tersebut
			currentCategoryData := categoryMap[inv.KategoriID]
			items := currentCategoryData["items"].([]gin.H)
			items = append(items, invData)
			currentCategoryData["items"] = items
			categoryMap[inv.KategoriID] = currentCategoryData
		}
	}

	// Konversi map ke array untuk response
	var result []gin.H
	for _, category := range categoryMap {
		result = append(result, category)
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
	
	// Cari inventaris yang akan dihapus untuk mengambil informasinya
	var inventaris types.Inventaris
	if err := db.Preload("SebaranBarang").First(&inventaris, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventaris tidak ditemukan"})
		return
	}

	// Log UserID for debugging
	fmt.Println("UserID:", inventaris.UserID)  // Add this line for debugging
	
	// Ambil data user
	var user types.User
	if inventaris.UserID == 0 {
		// Assign default user, e.g., user with ID 1 (or any default user you choose)
		if err := db.First(&user, 1).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user default"})
			return
		}
	} else {
		// Otherwise, fetch the user by UserID
		if err := db.First(&user, inventaris.UserID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data user"})
			return
		}
	}
	// Ambil data divisi
	var divisi types.Divisi
	if err := db.First(&divisi, inventaris.DivisiID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data divisi"})
		return
	}

	// Ambil posisi akhir dari SebaranBarang (jika ada)
	var posisiAkhir string
	if len(inventaris.SebaranBarang) > 0 {
		// Misalnya, posisi akhir ada pada field PosisiAkhir di SebaranBarang
		posisiAkhir = inventaris.SebaranBarang[len(inventaris.SebaranBarang)-1].PosisiAkhir
	} else {
		posisiAkhir = "Posisi tidak tersedia"
	}

	// Mulai transaksi
	tx := db.Begin()

	
	// Hapus Depresiasi yang terkait dengan Inventaris
	if err := tx.Where("id_barang = ?", id).Delete(&types.Depresiasi{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data depresiasi: " + err.Error()})
		return
	}

	// Hapus SebaranBarang yang terkait dengan Inventaris
	if err := tx.Where("id_barang = ?", id).Delete(&types.SebaranBarang{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data sebaran barang: " + err.Error()})
		return
	}
	
	// Buat catatan history sebelum menghapus
	now := time.Now()
	historyKeterangan := fmt.Sprintf("Barang %s yang ada di %s telah dijual oleh %s dari divisi %s sebanyak %d buah pada %s", 
		inventaris.NamaBarang,
		posisiAkhir, // Menambahkan posisi akhir dari SebaranBarang
		user.Name,
		divisi.NamaDivisi,
		inventaris.QtyBarang,
		now.Format("02-01-2006 15:04:05"))
	
	history := types.History{
		Kategori:   "Barang Keluar",
		Keterangan: historyKeterangan,
		CreatedAt:  now,
	}
	
	// Simpan data history ke database
	if err := tx.Create(&history).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data history: " + err.Error()})
		return
	}
	
	// Hapus inventaris
	if err := tx.Delete(&types.Inventaris{}, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// Commit transaksi jika semua operasi berhasil
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventaris berhasil dihapus",
		"history": history,
	})
}

func GetTotalInventaris(c *gin.Context) {
    db := database.GetDB()

    var totalQty int
    var totalTersedia int
    var totalTerpakai int
    var totalNilai int

    // Ambil semua data inventaris
    var inventaris []types.Inventaris
    if err := db.Find(&inventaris).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Hitung semua total
    for _, item := range inventaris {
        totalQty += item.QtyBarang
        totalTersedia += item.QtyTersedia
        totalTerpakai += item.QtyTerpakai
        totalNilai += item.TotalNilai
    }

    // Kirim hasil perhitungan sebagai JSON
    c.JSON(http.StatusOK, gin.H{
        "total_qty_barang":   totalQty,
        "total_qty_tersedia": totalTersedia,
        "total_qty_terpakai": totalTerpakai,
        "total_nilai_barang": totalNilai,
    })
}

