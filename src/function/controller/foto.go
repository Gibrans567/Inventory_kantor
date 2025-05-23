package controller

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"log"
	"os"

	"inventory/src/types"
	"inventory/src/function/database"

	"github.com/gin-gonic/gin"
)

func UploadGambarMulti(c *gin.Context) {
    // Mendapatkan ID dari query parameter
    id := c.DefaultQuery("id", "0")
    log.Printf("Search for Inventaris with ID: %s", id)

    // Mengambil data Inventaris berdasarkan ID
    inv := types.Inventaris{}
    db := database.GetDB() // Menggunakan GetDB untuk mendapatkan koneksi
    if err := db.First(&inv, id).Error; err != nil {
        log.Printf("Inventaris with ID %s not found: %v", id, err)
        c.JSON(http.StatusNotFound, gin.H{
            "status": "error",
            "message": "Inventaris not found",
            "foto": nil,
        })
        return
    }
    log.Printf("Inventaris found: ID %s", id)

    // Mendapatkan file gambar dari request
    form, err := c.MultipartForm()
    if err != nil {
        log.Printf("Error getting multipart form: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "error",
            "message": "Error processing uploaded files",
            "foto": nil,
        })
        return
    }
    
    files := form.File["link_foto"]
    if len(files) == 0 {
        log.Println("No file uploaded")
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "error",
            "message": "No file uploaded",
            "foto": nil,
        })
        return
    }
    log.Printf("Files uploaded: %d", len(files))

    // Validasi ekstensi file untuk memastikan hanya gambar yang diupload
    validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}

    var uploadPaths []string // Menyimpan path file yang berhasil di-upload

    // Proses setiap file
    for _, file := range files {
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
            c.JSON(http.StatusBadRequest, gin.H{
                "status": "error",
                "message": "Invalid file type. Only image files are allowed",
                "foto": nil,
            })
            return
        }

        // Menambahkan pembatasan ukuran file (maksimal 5MB)
        const MaxFileSize = 5 * 1024 * 1024 // 5 MB
        if file.Size > MaxFileSize {
            log.Println("File size exceeds the limit")
            c.JSON(http.StatusBadRequest, gin.H{
                "status": "error",
                "message": "File size exceeds the 5MB limit",
                "foto": nil,
            })
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
                c.JSON(http.StatusInternalServerError, gin.H{
                    "status": "error",
                    "message": "Failed to create storage directory",
                    "foto": nil,
                })
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
            c.JSON(http.StatusInternalServerError, gin.H{
                "status": "error",
                "message": "Failed to save file",
                "foto": nil,
            })
            return
        }
        log.Printf("File saved to: %s", filePath)

        // Membuat path yang bisa diakses melalui URL
        relativePath := strings.TrimPrefix(filepath.ToSlash(filePath), "./") // Menghilangkan './'
        uploadNotaPath := fmt.Sprintf("https://sysbar.awh.co.id/%s", relativePath)

        // Simpan path foto ke dalam tabel BarangFoto
        barangFoto := types.BarangFoto{
            IdBarang: inv.ID,     // ID barang yang terkait
            LinkFoto: uploadNotaPath,  // Path atau URL foto
        }

        // Simpan objek BarangFoto ke dalam database
        if err := db.Create(&barangFoto).Error; err != nil {
            log.Printf("Failed to save BarangFoto: %v", err)
            c.JSON(http.StatusInternalServerError, gin.H{
                "status": "error",
                "message": "Failed to save photo record",
                "foto": nil,
            })
            return
        }
        log.Println("BarangFoto record saved successfully")

        // Menyimpan path yang berhasil di-upload
        uploadPaths = append(uploadPaths, uploadNotaPath)
    }

    // Menampilkan hasil
    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "message": "Files uploaded successfully",
        "foto": uploadPaths,
    })
}

// GetBarangFoto retrieves a photo by ID
func GetBarangFoto(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "error",
            "message": "Invalid ID format",
            "foto": nil,
        })
        return
    }

    // Get the database connection
    db := database.GetDB()

    var barangFoto types.BarangFoto
    if err := db.First(&barangFoto, id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{
            "status": "error",
            "message": "Record not found",
            "foto": nil,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "message": "Photo found",
        "foto": barangFoto,
    })
}

// GetAllBarangFoto retrieves all photo records
func GetAllBarangFoto(c *gin.Context) {
    // Get the database connection
    db := database.GetDB()

    var barangFotoList []types.BarangFoto
    if err := db.Find(&barangFotoList).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "error",
            "message": err.Error(),
            "foto": nil,
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "message": "All photos retrieved",
        "foto": barangFotoList,
    })
}

func GetBarangFotoByBarang(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "error",
            "message": "Invalid ID format",
            "foto": nil,
        })
        return
    }

    // Get the database connection
    db := database.GetDB()

    var barangFotoList []types.BarangFoto
    if err := db.Where("id_barang = ?", id).Find(&barangFotoList).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "status": "error",
            "message": err.Error(),
            "foto": nil,
        })
        return
    }

    if len(barangFotoList) == 0 {
        c.JSON(http.StatusNotFound, gin.H{
            "status": "error",
            "message": "No photos found for this item",
            "foto": nil,
        })
        return
    }

    // Success response
    c.JSON(http.StatusOK, gin.H{
        "status": "success",
        "message": "Photos retrieved successfully",
        "foto": barangFotoList,
    })
}

// UpdateBarangFoto updates an existing photo record
func UpdateBarangFoto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get the database connection
	db := database.GetDB()

	// Check if record exists
	var existingFoto types.BarangFoto
	if err := db.First(&existingFoto, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	// Parse the form data
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to parse form"})
		return
	}

	// Update item ID if provided
	if idBarang, err := strconv.ParseUint(c.PostForm("id_barang"), 10, 32); err == nil {
		existingFoto.IdBarang = uint(idBarang)
	}

	// Update photo if provided
	file, fileHeader, err := c.Request.FormFile("foto")
	if err == nil {
		defer file.Close()

		// Check file type
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}
		
		// Cek apakah ekstensi file valid
		valid := false
		for _, e := range validExtensions {
			if ext == e {
				valid = true
				break
			}
		}

		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only image files are allowed"})
			return
		}

		// Menambahkan pembatasan ukuran file (maksimal 5MB)
		const MaxFileSize = 5 * 1024 * 1024 // 5 MB
		if fileHeader.Size > MaxFileSize {
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
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create storage directory"})
				return
			}
		}

		// Dapatkan informasi barang terkait untuk nama file
		var inv types.Inventaris
		if err := db.First(&inv, existingFoto.IdBarang).Error; err != nil {
			// Jika tidak bisa mendapatkan informasi inventaris, gunakan ID saja
			newFileName := fmt.Sprintf("barang_%d_%s%s", existingFoto.IdBarang, currentDate.Format("2006-01-02"), ext)
			filePath := filepath.Join(storageDir, newFileName)
			
			if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
				return
			}
			
			// Membuat path yang bisa diakses melalui URL
			relativePath := strings.TrimPrefix(filepath.ToSlash(filePath), "./") // Menghilangkan './'
			existingFoto.LinkFoto = fmt.Sprintf("https://sysbar.awh.co.id/%s", relativePath)
		} else {
			// Jika informasi inventaris tersedia, gunakan format yang sama dengan UploadGambarMulti
			tanggalPembelian := inv.TanggalPembelian.Format("2006-01-02")
			newFileName := fmt.Sprintf("%s_%s%s", inv.NamaBarang, tanggalPembelian, ext)
			filePath := filepath.Join(storageDir, newFileName)
			
			if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
				return
			}
			
			// Membuat path yang bisa diakses melalui URL
			relativePath := strings.TrimPrefix(filepath.ToSlash(filePath), "./") // Menghilangkan './'
			existingFoto.LinkFoto = fmt.Sprintf("https://sysbar.awh.co.id/%s", relativePath)
		}
	}

	// Save the updated record
	if err := db.Save(&existingFoto).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": existingFoto})
}

// DeleteBarangFoto deletes a photo record by ID
func DeleteBarangFoto(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get the database connection
	db := database.GetDB()

	// Delete the record
	result := db.Delete(&types.BarangFoto{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Record with ID %d not found", id)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Record deleted successfully"})
}

// DeleteAllBarangFotoByBarang deletes all photos for a specific item
func DeleteAllBarangFotoByBarang(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Get the database connection
	db := database.GetDB()

	// Delete all records with the specified barang ID
	result := db.Where("id_barang = ?", id).Delete(&types.BarangFoto{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All photos for this item deleted successfully"})
}

	
	