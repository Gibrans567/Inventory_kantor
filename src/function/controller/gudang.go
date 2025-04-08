package controller

import (
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateGudang - Menambahkan Gudang baru
func CreateGudang(c *gin.Context) {
	var gudang types.Gudang
	if err := c.ShouldBindJSON(&gudang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validation to check if required fields are provided
	if gudang.NamaGudang == "" || gudang.LokasiGudang == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NamaGudang and LokasiGudang are required"})
		return
	}

	db := database.GetDB()
	result := db.Create(&gudang)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gudang)
}



func GetAllGudang(c *gin.Context) {
    var gudangs []types.Gudang
    db := database.GetDB()

    result := db.Omit("Inventaris", "Depresiasi").
        Select("id", "nama_gudang", "lokasi_gudang", "created_at", "updated_at").
        Find(&gudangs)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch gudang data"})
        return
    }

    c.JSON(http.StatusOK, gudangs)
}




// UpdateGudang - Mengupdate data Gudang
func UpdateGudang(c *gin.Context) {
	id := c.Param("id")
	var gudang types.Gudang
	if err := c.ShouldBindJSON(&gudang); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()

	// Mengecek apakah Gudang dengan ID tersebut ada
	var existingGudang types.Gudang
	if err := db.First(&existingGudang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gudang not found"})
		return
	}

	// Pembaruan field yang diberikan
	result := db.Model(&existingGudang).Where("id = ?", id).Updates(gudang)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, existingGudang)
}


// DeleteGudang - Menghapus Gudang berdasarkan ID
func DeleteGudang(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	// Mengecek apakah Gudang dengan ID tersebut ada
	var gudang types.Gudang
	if err := db.First(&gudang, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gudang not found"})
		return
	}

	// Menghapus Gudang dengan ID yang diberikan
	result := db.Delete(&gudang)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Mengirimkan respons sukses setelah penghapusan
	c.JSON(http.StatusNoContent, gin.H{"message": "Deleted successfully"})
}
