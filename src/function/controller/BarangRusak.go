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

func CreateBarangStatus(ctx *gin.Context) {
	var barangStatus types.BarangStatus
	if err := ctx.ShouldBindJSON(&barangStatus); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	// Get the database connection
	db := database.GetDB()

	// Create the record
	if err := db.Create(&barangStatus).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusCreated, Response{
		Status:  "success",
		Message: "Barang status created successfully",
		Data:    barangStatus,
	})
}

func GetBarangStatus(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid ID format",
			Data:    nil,
		})
		return
	}

	// Get the database connection
	db := database.GetDB()

	var barangStatus types.BarangStatus
	if err := db.First(&barangStatus, id).Error; err != nil {
		ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Record not found",
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Record retrieved successfully",
		Data:    barangStatus,
	})
}

func GetAllBarangStatus(ctx *gin.Context) {
	// Get the database connection
	db := database.GetDB()

	var barangStatusList []types.BarangStatus
	if err := db.Find(&barangStatusList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "All records retrieved successfully",
		Data:    barangStatusList,
	})
}

func GetBarangStatusByBarang(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid ID format",
			Data:    nil,
		})
		return
	}

	// Get the database connection
	db := database.GetDB()

	var barangStatusList []types.BarangStatus
	if err := db.Where("id_barang = ?", id).Find(&barangStatusList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Records retrieved successfully",
		Data:    barangStatusList,
	})
}

func GetBarangStatusBySebaran(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
			Status:  "error",
			Message: "Invalid ID format",
			Data:    nil,
		})
		return
	}

	// Get the database connection
	db := database.GetDB()

	var barangStatusList []types.BarangStatus
	if err := db.Where("id_sebaran_barang = ?", id).Find(&barangStatusList).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Records retrieved successfully",
		Data:    barangStatusList,
	})
}

func UpdateBarangStatus(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
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
		ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: "Record not found",
			Data:    nil,
		})
		return
	}

	// Bind the JSON to the existing record
	if err := ctx.ShouldBindJSON(&existingStatus); err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
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
		ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Record updated successfully",
		Data:    existingStatus,
	})
}

func DeleteBarangStatus(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Response{
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
		ctx.JSON(http.StatusInternalServerError, Response{
			Status:  "error",
			Message: result.Error.Error(),
			Data:    nil,
		})
		return
	}

	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, Response{
			Status:  "error",
			Message: fmt.Sprintf("Record with ID %d not found", id),
			Data:    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, Response{
		Status:  "success",
		Message: "Record deleted successfully",
		Data:    nil,
	})
}
