package storage

import (
	"net/http"
	"inventory/src/types"
    "inventory/src/function/database"
    "time"

	"github.com/gin-gonic/gin"
)

func CreateInventory(c *gin.Context) {
    var input types.Inventory
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    input.TanggalBeli = time.Now()
    input.TanggalMutasi = time.Now()
    input.PosisiAkhir = nil // Pastikan posisi_akhir NULL saat insert

    db := database.GetDB()
    db.Create(&input)

    c.JSON(http.StatusOK, input)
}


func GetInventories(c *gin.Context) {
    var inventories []types.Inventory
    db := database.GetDB()
    db.Find(&inventories)
    c.JSON(http.StatusOK, inventories)
}

func GetInventoryByID(c *gin.Context) {
    var inventory types.Inventory
    db := database.GetDB()
    if err := db.First(&inventory, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
        return
    }
    c.JSON(http.StatusOK, inventory)
}

func UpdateInventory(c *gin.Context) {
    var inventory types.Inventory
    db := database.GetDB()
    
    if err := db.First(&inventory, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Data tidak ditemukan"})
        return
    }

    var input types.Inventory
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Validasi: PosisiAkhir tidak boleh kosong saat update
    if input.PosisiAkhir == nil || *input.PosisiAkhir == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Posisi Akhir harus diisi saat update"})
        return
    }

    db.Model(&inventory).Updates(input)
    c.JSON(http.StatusOK, inventory)
}


func DeleteInventory(c *gin.Context) {
    var inventory types.Inventory
    db := database.GetDB()
    if err := db.First(&inventory, c.Param("id")).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
        return
    }

    db.Delete(&inventory)
    c.JSON(http.StatusOK, gin.H{"message": "Data deleted"})
}
