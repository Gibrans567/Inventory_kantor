package routes

import (
	"inventory/src/function/storages"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Group route untuk inventory
	inventory := r.Group("/inventory")
	{
		inventory.POST("/", storage.CreateInventory)
		inventory.GET("/", storage.GetInventories)
		inventory.GET("/:id", storage.GetInventoryByID)
		inventory.PUT("/:id", storage.UpdateInventory)
		inventory.DELETE("/:id", storage.DeleteInventory)
	}

	return r
}
