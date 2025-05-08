package controller

import (
	"fmt"
	"time"
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"
    "log"
	"strings"



	"os"
	"path/filepath"
)


func ApplyDepresiasi(c *gin.Context) {
	// Ambil koneksi ke database
	db := database.GetDB()
	fmt.Println("Menjalankan proses depresiasi...")

	// Cari semua jadwal depresiasi yang memiliki tanggal yang sama dengan tanggal saat ini
	var jadwalDepresiasi []types.JadwalDepresiasi
	if err := db.Where("next_run <= ?", time.Now()).Find(&jadwalDepresiasi).Error; err != nil {
		fmt.Println("Gagal mendapatkan jadwal depresiasi:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan jadwal depresiasi"})
		return
	}

	// Jika tidak ada jadwal depresiasi yang ditemukan
	if len(jadwalDepresiasi) == 0 {
		fmt.Println("Tidak ada jadwal depresiasi yang sesuai dengan tanggal saat ini.")
		c.JSON(http.StatusOK, gin.H{"message": "Tidak ada jadwal depresiasi yang sesuai dengan tanggal saat ini."})
		return
	}

	// Proses setiap barang yang dijadwalkan untuk depresiasi
	for _, jadwal := range jadwalDepresiasi {
		// Cari inventaris berdasarkan barang yang ada pada jadwal
		var inventaris types.Inventaris
		if err := db.Where("id = ?", jadwal.IdBarang).First(&inventaris).Error; err != nil {
			fmt.Println("Gagal menemukan barang ID:", jadwal.IdBarang)
			continue // Jika barang tidak ditemukan, lanjut ke barang berikutnya
		}

		// Cari data depresiasi untuk barang ini
		var depresiasi types.Depresiasi
		if err := db.Where("id_barang = ?", jadwal.IdBarang).First(&depresiasi).Error; err != nil {
			fmt.Println("Gagal menemukan data depresiasi untuk barang ID:", jadwal.IdBarang)
			continue // Jika data depresiasi tidak ditemukan, lanjut ke barang berikutnya
		}

		// Hitung jumlah depresiasi berdasarkan nilai per bulan
		nilaiDepresiasi := depresiasi.Perbulan
		persentase := 2.5 // Anda bisa menyesuaikan persentase ini sesuai kebutuhan
		
		// Update hargaPembelian dengan mengurangi jumlah depresiasi
		newHargaPembelian := inventaris.HargaPembelian - nilaiDepresiasi
		
		// Hitung nilai total baru berdasarkan jumlah barang
		newTotalNilai := newHargaPembelian * inventaris.QtyBarang

		// Simpan perubahan harga pembelian dan total nilai ke dalam database
		if err := db.Model(&inventaris).Updates(map[string]interface{}{
			"harga_pembelian": newHargaPembelian,
			"total_nilai": newTotalNilai,
		}).Error; err != nil {
			fmt.Println("Gagal menyimpan nilai baru setelah depresiasi untuk barang ID:", jadwal.IdBarang)
			continue // Jika gagal menyimpan, lanjut ke barang berikutnya
		}

		// Buat catatan riwayat depresiasi
		now := time.Now()
		history := types.History{
			Kategori:   "Depresiasi",
			Keterangan: fmt.Sprintf("Harga %s telah dikurangi %.2f%% dari harga pembelian yang asli, yaitu %.2d",inventaris.NamaBarang, persentase, nilaiDepresiasi),
			CreatedAt:  now,
		}

		if err := db.Create(&history).Error; err != nil {
			fmt.Println("Gagal mencatat history depresiasi untuk barang ID:", jadwal.IdBarang)
		}

		// Menampilkan hasil depresiasi untuk debugging
		fmt.Printf("Depresiasi untuk barang ID %d berhasil diterapkan. Nilai baru: HargaPembelian=%.2d, TotalNilai=%.2d\n", 
			jadwal.IdBarang, newHargaPembelian, newTotalNilai)

            newNextRun := time.Now().AddDate(0, 1, 0) // Menambahkan 1 bulan ke depan

            // Update jadwal depresiasi dengan tanggal baru
            if err := db.Model(&jadwal).Update("next_run", newNextRun).Error; err != nil {
                fmt.Println("Gagal memperbarui jadwal depresiasi untuk barang ID:", jadwal.IdBarang)
            } else {
                // Menampilkan hasil update jadwal untuk debugging
                fmt.Printf("Jadwal depresiasi untuk barang ID %d telah diperbarui ke tanggal: %v\n", jadwal.IdBarang, newNextRun)
            }
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Depresiasi untuk barang yang dijadwalkan berhasil diterapkan.",
	})
}

func GetAllJadwal(c *gin.Context) {
	var histories []types.JadwalDepresiasi
	db := database.GetDB()
	result := db.Find(&histories)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch histories"})
		return
	}

	c.JSON(http.StatusOK, histories)
}

func RunDepresiationScheduler(interval time.Duration) {
	go func() {
		for {
			log.Println("Running depresiasi scheduler check...")
			processScheduledDepresiasi()
			
			// Tunggu sampai interval berikutnya
			time.Sleep(interval)
		}
	}()
	
	log.Printf("Depresiasi scheduler started with interval: %v\n", interval)
}

// processScheduledDepresiasi memeriksa jadwal depresiasi dan menjalankan
// proses depresiasi untuk jadwal yang waktunya telah tiba
func processScheduledDepresiasi() {
	// Ambil koneksi ke database
	db := database.GetDB()
	fmt.Println("Menjalankan proses depresiasi...")

	// Cari semua jadwal depresiasi yang memiliki tanggal yang sama dengan tanggal saat ini
	var jadwalDepresiasi []types.JadwalDepresiasi
	if err := db.Where("next_run <= ?", time.Now()).Find(&jadwalDepresiasi).Error; err != nil {
		fmt.Println("Gagal mendapatkan jadwal depresiasi:", err)
		return
	}

	// Jika tidak ada jadwal depresiasi yang ditemukan
	if len(jadwalDepresiasi) == 0 {
		fmt.Println("Tidak ada jadwal depresiasi yang sesuai dengan tanggal saat ini.")
		return
	}

	// Proses setiap barang yang dijadwalkan untuk depresiasi
	for _, jadwal := range jadwalDepresiasi {
		// Cari inventaris berdasarkan barang yang ada pada jadwal
		var inventaris types.Inventaris
		if err := db.Where("id = ?", jadwal.IdBarang).First(&inventaris).Error; err != nil {
			fmt.Println("Gagal menemukan barang ID:", jadwal.IdBarang)
			continue // Jika barang tidak ditemukan, lanjut ke barang berikutnya
		}

		// Cari data depresiasi untuk barang ini
		var depresiasi types.Depresiasi
		if err := db.Where("id_barang = ?", jadwal.IdBarang).First(&depresiasi).Error; err != nil {
			fmt.Println("Gagal menemukan data depresiasi untuk barang ID:", jadwal.IdBarang)
			continue // Jika data depresiasi tidak ditemukan, lanjut ke barang berikutnya
		}

		// Hitung jumlah depresiasi berdasarkan nilai per bulan
		nilaiDepresiasi := depresiasi.Perbulan
		persentase := 2.5 // Anda bisa menyesuaikan persentase ini sesuai kebutuhan
		
		// Update hargaPembelian dengan mengurangi jumlah depresiasi
		newHargaPembelian := inventaris.HargaPembelian - nilaiDepresiasi
		
		// Hitung nilai total baru berdasarkan jumlah barang
		newTotalNilai := newHargaPembelian * inventaris.QtyBarang

		// Simpan perubahan harga pembelian dan total nilai ke dalam database
		if err := db.Model(&inventaris).Updates(map[string]interface{}{
			"harga_pembelian": newHargaPembelian,
			"total_nilai": newTotalNilai,
		}).Error; err != nil {
			fmt.Println("Gagal menyimpan nilai baru setelah depresiasi untuk barang ID:", jadwal.IdBarang)
			continue // Jika gagal menyimpan, lanjut ke barang berikutnya
		}

		// Buat catatan riwayat depresiasi
		now := time.Now()
		history := types.History{
			Kategori:   "Depresiasi",
			Keterangan: fmt.Sprintf("Harga %s telah dikurangi %.2f%% dari harga pembelian yang asli, yaitu %d",
				inventaris.NamaBarang, persentase, nilaiDepresiasi),
			CreatedAt:  now,
		}

		if err := db.Create(&history).Error; err != nil {
			fmt.Println("Gagal mencatat history depresiasi untuk barang ID:", jadwal.IdBarang)
		}

		// Menampilkan hasil depresiasi untuk debugging
		fmt.Printf("Depresiasi untuk barang ID %d berhasil diterapkan. Nilai baru: HargaPembelian=%d, TotalNilai=%d\n", 
			jadwal.IdBarang, newHargaPembelian, newTotalNilai)

		newNextRun := time.Now().AddDate(0, 1, 0) // Menambahkan 1 bulan ke depan

		// Update jadwal depresiasi dengan tanggal baru
		if err := db.Model(&jadwal).Update("next_run", newNextRun).Error; err != nil {
			fmt.Println("Gagal memperbarui jadwal depresiasi untuk barang ID:", jadwal.IdBarang)
		} else {
			// Menampilkan hasil update jadwal untuk debugging
			fmt.Printf("Jadwal depresiasi untuk barang ID %d telah diperbarui ke tanggal: %v\n", jadwal.IdBarang, newNextRun)
		}
	}

	fmt.Println("Proses depresiasi otomatis selesai.")
}

// InitiateScheduler memulai scheduler saat aplikasi dijalankan
func InitiateScheduler() {
	// Jalankan scheduler setiap hari pada jam 00:00
	// RunDepresiationScheduler(24 * time.Hour)
    RunDepresiationScheduler(1 * time.Hour)
}

// DeleteAllByTimeframe menghapus semua data pada rentang waktu tertentu dari semua tabel
func DeleteAllByTimeframe(c *gin.Context) {
    // Ambil data langsung dari request body
    var requestBody map[string]string
    if err := c.ShouldBindJSON(&requestBody); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    startDate, err := time.Parse("2006-01-02", requestBody["start_date"])
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal mulai tidak valid. Gunakan format YYYY-MM-DD"})
        return
    }

    endDate, err := time.Parse("2006-01-02", requestBody["end_date"])
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal akhir tidak valid. Gunakan format YYYY-MM-DD"})
        return
    }

    endDate = endDate.Add(24 * time.Hour)

    db := database.GetDB()
    deletedCounts := make(map[string]int64)

    // Daftar tabel dan field waktu masing-masing
    tables := []struct {
        name      string
        model     interface{}
        timeField string // Kolom yang digunakan untuk filter waktu
    }{
        {"sebaran_barang", &types.SebaranBarang{}, "created_at"},
        {"depresiasi", &types.Depresiasi{}, "created_at"},
        // JadwalDepresiasi sepertinya tidak memiliki kolom created_at, gunakan next_run
        {"jadwal_depresiasi", &types.JadwalDepresiasi{}, "next_run"},
        {"inventaris", &types.Inventaris{}, "created_at"},
        {"user", &types.User{}, "created_at"},
        {"divisi", &types.Divisi{}, "created_at"},
        {"kategori", &types.Kategori{}, "created_at"},
        {"gudang", &types.Gudang{}, "created_at"},
        {"history", &types.History{}, "created_at"},
    }

    tx := db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    for _, table := range tables {
		// Gunakan field waktu yang sesuai untuk masing-masing tabel
		deleteResult := tx.Where(table.timeField+" BETWEEN ? AND ?", startDate, endDate).Delete(table.model)
		
		if deleteResult.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Gagal menghapus data di tabel " + table.name,
				"message": deleteResult.Error.Error(),
			})
			return
		}
		
		deletedCounts[table.name] = deleteResult.RowsAffected
	}

    // Tambahkan ke history
    history := types.History{
        Kategori:   "DELETE",
        Keterangan: "Menghapus semua data dari " + requestBody["start_date"] + " hingga " + requestBody["end_date"],
    }
    
    if err := tx.Create(&history).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Gagal mencatat history",
            "message": err.Error(),
        })
        return
    }

    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Gagal menyelesaikan transaksi",
            "message": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message": "Semua data dalam rentang waktu berhasil dihapus",
        "deleted": deletedCounts,
    })
}

func UploadGambar(c *gin.Context) {
    // Mendapatkan ID dari query parameter
    id := c.DefaultQuery("id", "0")
    log.Printf("Search for Inventaris with ID: %s", id)

    // Mengambil data Inventaris berdasarkan ID
    inv := types.Inventaris{}
    db := database.GetDB() // Menggunakan GetDB untuk mendapatkan koneksi
    if err := db.First(&inv, id).Error; err != nil {
        log.Printf("Inventaris with ID %s not found: %v", id, err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Inventaris not found"})
        return
    }
    log.Printf("Inventaris found: ID %s", id)

    // Mendapatkan file gambar dari request
    file, _ := c.FormFile("upload_nota")
    if file == nil {
        log.Println("No file uploaded")
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
        return
    }
    log.Printf("File uploaded: %s", file.Filename)

    // Validasi ekstensi file untuk memastikan hanya gambar yang diupload
    validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}
    ext := strings.ToLower(filepath.Ext(file.Filename))

    // Cek apakah ekstensi file valid
    valid := false
    for _, e := range validExtensions {
        if ext == e {
            valid = true
            break
        }
    }

    if !valid {
        log.Println("Invalid file type. Only image files are allowed")
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only image files are allowed"})
        return
    }

    // Menambahkan pembatasan ukuran file (maksimal 5MB)
    const MaxFileSize = 5 * 1024 * 1024 // 5 MB
    if file.Size > MaxFileSize {
        log.Println("File size exceeds the limit")
        c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds the 5MB limit"})
        return
    }

    // Membuat folder penyimpanan dengan format: storage/tahun/bulan/hari
    currentDate := time.Now()
    storageDir := fmt.Sprintf("./storage/%d/%02d/%02d", currentDate.Year(), currentDate.Month(), currentDate.Day())

    // Cek jika folder belum ada, buat foldernya
    if _, err := os.Stat(storageDir); os.IsNotExist(err) {
        err := os.MkdirAll(storageDir, os.ModePerm)
        if err != nil {
            log.Printf("Failed to create storage directory: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create storage directory"})
            return
        }
        log.Println("Storage directory created successfully")
    }

    // Membuat nama file berdasarkan nama_barang dan tanggal_pembelian
    tanggalPembelian := inv.TanggalPembelian.Format("2006-01-02")
    newFileName := fmt.Sprintf("%s_%s%s", inv.NamaBarang, tanggalPembelian, ext)
    log.Printf("Generated new file name: %s", newFileName)

    // Menyimpan file ke folder storage
    filePath := filepath.Join(storageDir, newFileName)
    if err := c.SaveUploadedFile(file, filePath); err != nil {
        log.Printf("Failed to save file: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }
    log.Printf("File saved to: %s", filePath)

    // Membuat path yang bisa diakses melalui URL
    relativePath := strings.TrimPrefix(filepath.ToSlash(filePath), "./") // Menghilangkan './'
    uploadNotaPath := fmt.Sprintf("http://localhost:8080/%s", relativePath)

    // Update record Inventaris dengan path file
    inv.UploadNota = uploadNotaPath
    if err := db.Save(&inv).Error; err != nil {
        log.Printf("Failed to update database: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database"})
        return
    }
    log.Println("Database updated with new file path")

    // Menampilkan hasil
    c.JSON(http.StatusOK, gin.H{
        "message": "File uploaded successfully",
        "file":    uploadNotaPath,
    })
}

