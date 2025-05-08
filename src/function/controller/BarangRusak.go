package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"inventory/src/types"
	"inventory/src/function/database"

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

	// Set posisi terakhir dari sebaran
	barangStatus.PosisiAkhir = sebaran.PosisiAkhir

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
		inventaris.QtyTersedia -= barangStatus.QtyBarang

		if err := db.Save(&sebaran).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Gagal update sebaran",
				Data:    nil,
			})
			return
		}

		if err := db.Save(&inventaris).Error; err != nil {
			c.JSON(http.StatusInternalServerError, Response{
				Status:  "error",
				Message: "Gagal update inventaris",
				Data:    nil,
			})
			return
		}
	}

	if err := db.Create(&barangStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid ID format",
			Data:    nil,
		})
		return
	}

	// Get the database connection
	db := database.GetDB()

	// Check if record exists
	var existingStatus types.BarangStatus
	if err := db.First(&existingStatus, id).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Record not found",
			Data:    nil,
		})
		return
	}

	// Bind the JSON to the existing record
	if err := c.ShouldBindJSON(&existingStatus); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Ensure ID remains unchanged
	existingStatus.ID = uint(id)

	// Save the updated record
	if err := db.Save(&existingStatus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Record updated successfully",
		Data:    existingStatus,
	})
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

	// Get the database connection
	db := database.GetDB()

	// Delete the record
	result := db.Delete(&types.BarangStatus{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: result.Error.Error(),
			Data:    nil,
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: fmt.Sprintf("Record with ID %d not found", id),
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Record deleted successfully",
		Data:    nil,
	})
}
