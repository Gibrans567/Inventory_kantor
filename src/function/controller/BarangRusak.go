package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"inventory/src/types"
	"inventory/src/function/database"
	"time"


	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// Response structure
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"dataBarang"`
}

func CreateBarangStatus(c *gin.Context) {
	var barangStatus types.BarangStatus

	if err := c.ShouldBindJSON(&barangStatus); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	db := database.GetDB()

	var sebaran types.SebaranBarang
	if err := db.First(&sebaran, barangStatus.IdSebaranBarang).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Sebaran barang tidak ditemukan",
			Data:    nil,
		})
		return
	}

	var inventaris types.Inventaris
	if err := db.First(&inventaris, barangStatus.IdBarang).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Inventaris tidak ditemukan",
			Data:    nil,
		})
		return
	}

	barangStatus.PosisiAkhir = sebaran.PosisiAkhir

	var existingBarangStatus types.BarangStatus
	result := db.Where("id_barang = ? AND status = ? AND note = ?",
		barangStatus.IdBarang, barangStatus.Status, barangStatus.Note).First(&existingBarangStatus)

	var user types.User
	if err := db.First(&user, sebaran.IdUser).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "User tidak ditemukan",
			Data:    nil,
		})
		return
	}

	var divisi types.Divisi
	if err := db.First(&divisi, user.IdDivisi).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Divisi tidak ditemukan",
			Data:    nil,
		})
		return
	}

	now := time.Now()

	// Jika barang rusak atau maintenance, validasi qty
	if barangStatus.Status == "Barang rusak" || barangStatus.Status == "Maintenance" {
		if barangStatus.QtyBarang <= 0 {
			c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Qty barang harus lebih dari 0",
				Data:    nil,
			})
			return
		}

		if sebaran.QtyBarang < barangStatus.QtyBarang {
			c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Jumlah barang melebihi stok di sebaran",
				Data:    nil,
			})
			return
		}

		sebaran.QtyBarang -= barangStatus.QtyBarang
		inventaris.QtyRusak += barangStatus.QtyBarang

		if err := db.Save(&sebaran).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Gagal update sebaran",
				Data:    nil,
			})
			return
		}
	}

	// Jika record sudah ada, update
	if result.Error == nil {
		oldQty := existingBarangStatus.QtyBarang
		newQty := oldQty + barangStatus.QtyBarang
		existingBarangStatus.QtyBarang = newQty

		if err := db.Save(&existingBarangStatus).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Gagal update barang status",
				Data:    nil,
			})
			return
		}

		if inventaris.QtyTerpakai < barangStatus.QtyBarang {
			c.JSON(http.StatusBadRequest, Response{
				Status:  "error",
				Message: "Qty terpakai lebih kecil dari jumlah yang dikurangi",
				Data:    nil,
			})
			return
		}

		inventaris.QtyTerpakai -= barangStatus.QtyBarang
		if err := db.Save(&inventaris).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Gagal update qty_terpakai di inventaris",
				Data:    nil,
			})
			return
		}

		if barangStatus.Status == "Barang rusak" {
			history := types.History{
				Kategori: "Barang Rusak",
				Keterangan: fmt.Sprintf(
					"Barang %s telah dinyatakan rusak sebanyak %d oleh %s dari %s",
					inventaris.NamaBarang,
					barangStatus.QtyBarang,
					user.Name,
					divisi.NamaDivisi,
				),
				CreatedAt: now,
			}
			if err := db.Create(&history).Error; err != nil {
				c.JSON(http.StatusInternalServerError, Response{
					Status:  "error",
					Message: "Gagal menyimpan history",
					Data:    nil,
				})
				return
			}
		}

		c.JSON(http.StatusOK, Response{
			Status:  "success",
			Message: "Barang status updated successfully",
			Data:    existingBarangStatus,
		})
		return
	}

	// Jika record belum ada
	if err := db.Create(&barangStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if inventaris.QtyTerpakai < barangStatus.QtyBarang {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Qty terpakai lebih kecil dari jumlah yang dikurangi",
			Data:    nil,
		})
		return
	}

	inventaris.QtyTerpakai -= barangStatus.QtyBarang
	if err := db.Save(&inventaris).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Gagal update qty_terpakai di inventaris",
			Data:    nil,
		})
		return
	}

	if barangStatus.Status == "Barang rusak" {
		history := types.History{
			Kategori: "Barang Rusak",
			Keterangan: fmt.Sprintf(
				"Barang %s telah dinyatakan rusak sebanyak %d oleh %s dari %s",
				inventaris.NamaBarang,
				barangStatus.QtyBarang,
				user.Name,
				divisi.NamaDivisi,
			),
			CreatedAt: now,
		}
		if err := db.Create(&history).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Gagal menyimpan history",
				Data:    nil,
			})
			return
		}
	}

	c.JSON(http.StatusCreated, Response{
		Status:  "success",
		Message: "Barang status created successfully",
		Data:    barangStatus,
	})
}


func GetBarangStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid ID format",
			Data:    nil,
		})
		return
	}

	db := database.GetDB()
	var barangStatus types.BarangStatus

	if err := db.Preload("Inventaris").First(&barangStatus, id).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Record not found",
			Data:    nil,
		})
		return
	}

	responseData := map[string]interface{}{
		"id":           barangStatus.ID,
		"nama_barang":  barangStatus.Inventaris.NamaBarang,
		"status":       barangStatus.Status,
		"note":         barangStatus.Note,
		"qty_barang":   barangStatus.QtyBarang,
	}
	

	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Record retrieved successfully",
		Data:    responseData,
	})
}


func GetAllBarangStatus(c *gin.Context) {
	db := database.GetDB()

	var barangStatuses []types.BarangStatus
	if err := db.Preload("SebaranBarang.Inventaris").Find(&barangStatuses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	var responseData []map[string]interface{}
	for _, bs := range barangStatuses {
		responseData = append(responseData, map[string]interface{}{
			"id":           bs.ID,
			"nama_barang":  bs.SebaranBarang.Inventaris.NamaBarang,
			"status":       bs.Status,
			"note":         bs.Note,
			"qty_barang":   bs.QtyBarang,
		})
		
	}

	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Barang status records retrieved successfully",
		Data:    responseData,
	})
}


func GetBarangStatusByBarang(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid ID format",
			Data:    nil,
		})
		return
	}

	db := database.GetDB()
	var barangStatusList []types.BarangStatus

	if err := db.Preload("Inventaris").Where("id_barang = ?", id).Find(&barangStatusList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	var responseData []map[string]interface{}
	for _, bs := range barangStatusList {
		responseData = append(responseData, map[string]interface{}{
			"id":           bs.ID,
			"id_barang":   bs.IdBarang,
			"nama_barang":  bs.Inventaris.NamaBarang,
			"status":       bs.Status,
			"note":         bs.Note,
			"qty_barang":   bs.QtyBarang,
		})
		
	}

	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Records retrieved successfully",
		Data:    responseData,
	})
}


func GetBarangStatusBySebaran(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid ID format",
			Data:    nil,
		})
		return
	}

	db := database.GetDB()
	var barangStatusList []types.BarangStatus

	if err := db.Preload("Inventaris").Where("id_sebaran_barang = ?", id).Find(&barangStatusList).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	var responseData []map[string]interface{}
	for _, bs := range barangStatusList {
		responseData = append(responseData, map[string]interface{}{
			"id":           bs.ID,
			"nama_barang":  bs.Inventaris.NamaBarang,
			"status":       bs.Status,
			"note":         bs.Note,
			"qty_barang":   bs.QtyBarang,
		})
		
	}

	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Records retrieved successfully",
		Data:    responseData,
	})
}


func UpdateBarangStatus(c *gin.Context) {
	id := c.Param("id")

	barangStatusID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid ID format",
			Data:    nil,
		})
		return
	}

	db := database.GetDB()

	var existingBarangStatus types.BarangStatus
	if err := db.First(&existingBarangStatus, barangStatusID).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Barang status tidak ditemukan",
			Data:    nil,
		})
		return
	}

	var updatedBarangStatus types.BarangStatus
	if err := c.ShouldBindJSON(&updatedBarangStatus); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	if updatedBarangStatus.QtyBarang > existingBarangStatus.QtyBarang {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: fmt.Sprintf("Qty Barang tidak boleh lebih dari %d", existingBarangStatus.QtyBarang),
			Data:    nil,
		})
		return
	}

	var inventaris types.Inventaris
	if err := db.First(&inventaris, existingBarangStatus.IdBarang).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Inventaris tidak ditemukan",
			Data:    nil,
		})
		return
	}

	oldStatus := existingBarangStatus.Status

	isStatusBecomingAvailable := updatedBarangStatus.Status == "Tersedia" ||
		updatedBarangStatus.Status == "Sudah fix" ||
		updatedBarangStatus.Status == "Sudah bisa digunakan" ||
		updatedBarangStatus.Status == "Clear"

	wasUnavailableStatus := oldStatus == "Barang rusak" || oldStatus == "Maintenance"

	if wasUnavailableStatus && isStatusBecomingAvailable {
		inventaris.QtyTersedia += updatedBarangStatus.QtyBarang

		if inventaris.QtyRusak >= updatedBarangStatus.QtyBarang {
			inventaris.QtyRusak -= updatedBarangStatus.QtyBarang
		} else {
			inventaris.QtyRusak = 0
		}

		if err := db.Save(&inventaris).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Gagal update inventaris",
				Data:    nil,
			})
			return
		}

		if updatedBarangStatus.QtyBarang < existingBarangStatus.QtyBarang {
			remainingQty := existingBarangStatus.QtyBarang - updatedBarangStatus.QtyBarang
			existingBarangStatus.QtyBarang = remainingQty

			if err := db.Save(&existingBarangStatus).Error; err != nil {
				c.JSON(http.StatusInternalServerError, Response{
					Status:  "error",
					Message: "Gagal update barang status",
					Data:    nil,
				})
				return
			}

			c.JSON(http.StatusOK, Response{
				Status:  "success",
				Message: fmt.Sprintf("%d item berhasil diubah menjadi tersedia, %d item masih dengan status %s",
					updatedBarangStatus.QtyBarang, remainingQty, oldStatus),
				Data: existingBarangStatus,
			})
			return
		} else {
			// History: Catat siapa yang memperbaiki atau service
			var sebaran types.SebaranBarang
			if err := db.First(&sebaran, existingBarangStatus.IdSebaranBarang).Error; err == nil {
				var user types.User
				if err := db.First(&user, sebaran.IdUser).Error; err == nil {
					var divisi types.Divisi
					if err := db.First(&divisi, user.IdDivisi).Error; err == nil {
						var history types.History
						if oldStatus == "Barang rusak" {
							history = types.History{
								Kategori:   "Barang Diperbaiki",
								Keterangan: fmt.Sprintf("Barang %s telah diperbaiki dan siap digunakan oleh %s dari %s", inventaris.NamaBarang, user.Name, divisi.NamaDivisi),
								CreatedAt:  time.Now(),
							}
						} else if oldStatus == "Maintenance" {
							history = types.History{
								Kategori:   "Barang Selesai Service",
								Keterangan: fmt.Sprintf("Barang %s telah selesai di service oleh %s dari %s", inventaris.NamaBarang, user.Name, divisi.NamaDivisi),
								CreatedAt:  time.Now(),
							}
						}
						if history.Kategori != "" {
							_ = db.Create(&history)
						}
					}
				}
			}

			// Hapus existingBarangStatus karena semua qty sudah tersedia
			if err := db.Delete(&existingBarangStatus).Error; err != nil {
				c.JSON(http.StatusInternalServerError, Response{
					Status:  "error",
					Message: "Gagal menghapus barang status",
					Data:    nil,
				})
				return
			}

			// **Cek dan hapus sebaranBarang jika QtyBarang <= 0**
			if err := db.First(&sebaran, existingBarangStatus.IdSebaranBarang).Error; err == nil {
				if sebaran.QtyBarang <= 0 {
					if err := db.Delete(&sebaran).Error; err != nil {
						c.JSON(http.StatusInternalServerError, Response{
							Status:  "error",
							Message: "Gagal menghapus data sebaran barang",
							Data:    nil,
						})
						return
					}
				}
			}

			c.JSON(http.StatusOK, Response{
				Status:  "success",
				Message: "Semua barang berhasil diubah menjadi tersedia dan record dihapus",
				Data:    existingBarangStatus,
			})
			return
		}
	}

	if isStatusBecomingAvailable && updatedBarangStatus.QtyBarang != existingBarangStatus.QtyBarang {
		qtyDiff := updatedBarangStatus.QtyBarang - existingBarangStatus.QtyBarang
		inventaris.QtyTersedia += qtyDiff

		if inventaris.QtyRusak >= updatedBarangStatus.QtyBarang {
			inventaris.QtyRusak -= updatedBarangStatus.QtyBarang
		} else {
			inventaris.QtyRusak = 0
		}

		if err := db.Save(&inventaris).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Gagal update inventaris",
				Data:    nil,
			})
			return
		}

		if updatedBarangStatus.QtyBarang < existingBarangStatus.QtyBarang {
			remainingQty := existingBarangStatus.QtyBarang - updatedBarangStatus.QtyBarang
			existingBarangStatus.QtyBarang = remainingQty

			if err := db.Save(&existingBarangStatus).Error; err != nil {
				c.JSON(http.StatusInternalServerError, Response{
					Status:  "error",
					Message: "Gagal update barang status",
					Data:    nil,
				})
				return
			}

			c.JSON(http.StatusOK, Response{
				Status:  "success",
				Message: fmt.Sprintf("%d item berhasil diubah menjadi tersedia, %d item masih dengan status %s",
					updatedBarangStatus.QtyBarang, remainingQty, oldStatus),
				Data: existingBarangStatus,
			})
			return
		} else {
			if err := db.Delete(&existingBarangStatus).Error; err != nil {
				c.JSON(http.StatusInternalServerError, Response{
					Status:  "error",
					Message: "Gagal menghapus barang status",
					Data:    nil,
				})
				return
			}

			// **Cek dan hapus sebaranBarang jika QtyBarang <= 0**
			var sebaran types.SebaranBarang
			if err := db.First(&sebaran, existingBarangStatus.IdSebaranBarang).Error; err == nil {
				if sebaran.QtyBarang <= 0 {
					if err := db.Delete(&sebaran).Error; err != nil {
						c.JSON(http.StatusInternalServerError, Response{
							Status:  "error",
							Message: "Gagal menghapus data sebaran barang",
							Data:    nil,
						})
						return
					}
				}
			}

			c.JSON(http.StatusOK, Response{
				Status:  "success",
				Message: "Semua barang berhasil diubah menjadi tersedia dan record dihapus",
				Data:    existingBarangStatus,
			})
			return
		}
	}

	if !isStatusBecomingAvailable {
		existingBarangStatus.Status = updatedBarangStatus.Status
		existingBarangStatus.QtyBarang = updatedBarangStatus.QtyBarang

		if updatedBarangStatus.Note != "" {
			existingBarangStatus.Note = updatedBarangStatus.Note
		}
		if updatedBarangStatus.PosisiAkhir != "" {
			existingBarangStatus.PosisiAkhir = updatedBarangStatus.PosisiAkhir
		}

		if err := db.Save(&existingBarangStatus).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Gagal update barang status",
				Data:    nil,
			})
			return
		}

		c.JSON(http.StatusOK, Response{
			Status:  "success",
			Message: "Barang status updated successfully",
			Data:    existingBarangStatus,
		})
	}
}


func DeleteBarangStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid ID format",
			Data:    nil,
		})
		return
	}

	db := database.GetDB()

	// Ambil data BarangStatus terlebih dahulu
	var barangStatus types.BarangStatus
	if err := db.First(&barangStatus, id).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: fmt.Sprintf("Record with ID %d not found", id),
			Data:    nil,
		})
		return
	}

	qtyToDeduct := barangStatus.QtyBarang
	inventarisID := barangStatus.IdBarang // Sesuai dengan relasi di struct

	// Hapus data BarangStatus
	result := db.Delete(&barangStatus)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: result.Error.Error(),
			Data:    nil,
		})
		return
	}

	// Kurangi qty_barang di tabel inventaris
	if err := db.Model(&types.Inventaris{}).
		Where("id = ?", inventarisID).
		Update("qty_barang", gorm.Expr("qty_barang - ?", qtyToDeduct)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: "Record deleted, but failed to update inventaris: " + err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Record deleted and inventaris updated successfully",
		Data:    nil,
	})
}

