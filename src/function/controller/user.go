package controller

import (
	"inventory/src/types"
	"inventory/src/function/database"
	"github.com/gin-gonic/gin"
	"net/http"


)

// CreateUser - Menambahkan User baru
func CreateUser(c *gin.Context) {
	var user types.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	result := db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUserByID - Mendapatkan User berdasarkan ID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	var user types.User
	db := database.GetDB()
	result := db.First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetAllUsers(c *gin.Context) {
    var users []types.User
    db := database.GetDB()

    // Ambil semua user dari database
    result := db.Find(&users)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        return
    }

    // Kembalikan response dengan status OK dan data users
    c.JSON(http.StatusOK, users)
}

// UpdateUser - Mengupdate data User
func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user types.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := database.GetDB()
	result := db.Model(&types.User{}).Where("id = ?", id).Updates(user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser - Menghapus User berdasarkan ID
func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	db := database.GetDB()

	result := db.Delete(&types.User{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{"message": "Deleted successfully"})
}

