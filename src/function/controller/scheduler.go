package controller

import (
	"fmt"
	"time"
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"
    "log"
)


func CreateInventaris1(c *gin.Context) {
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
        Keterangan: fmt.Sprintf("Pada %s barang %s telah masuk ke gudang %s", now.Format("02-01-2006 15:04:05"), inventaris.NamaBarang, gudang.NamaGudang),
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