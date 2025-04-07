package routes

import (
	"inventory/src/function/storages"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Depresiasi Routes
	depresiasi := r.Group("/depresiasi")
	{
		depresiasi.POST("/", storage.CreateDepresiasi)
		depresiasi.GET("/", storage.GetAllDepresiasi)
		depresiasi.PUT("/:id", storage.UpdateDepresiasi)
		depresiasi.DELETE("/:id", storage.DeleteDepresiasi)
	}

	// Divisi Routes
	divisi := r.Group("/divisi")
	{
		divisi.POST("/", storage.CreateDivisi)
		divisi.GET("/", storage.GetAllDivisi)
		divisi.PUT("/:id", storage.UpdateDivisi)
		divisi.DELETE("/:id", storage.DeleteDivisi)
	}

	// Gudang Routes
	gudang := r.Group("/gudang")
	{
		gudang.POST("/", storage.CreateGudang)
		gudang.GET("/", storage.GetAllGudang)
		gudang.POST("/update/:id", storage.UpdateGudang)
		gudang.DELETE("/:id", storage.DeleteGudang)
	}

	// Inventaris Routes
	inventaris := r.Group("/inventaris")
	{
		inventaris.POST("/", storage.CreateInventaris)
		inventaris.GET("/", storage.GetAllInventaris)
		inventaris.PUT("/:id", storage.UpdateInventaris)
		inventaris.DELETE("/:id", storage.DeleteInventaris)
	}

	// Kategori Routes
	kategori := r.Group("/kategori")
	{
		kategori.POST("/", storage.CreateKategori)
		kategori.GET("/", storage.GetAllKategori)
		kategori.PUT("/:id", storage.UpdateKategori)
		kategori.DELETE("/:id", storage.DeleteKategori)
	}

	// SebaranBarang Routes
	sebaranBarang := r.Group("/sebaranBarang")
	{
		sebaranBarang.POST("/", storage.CreateSebaranBarang)
		sebaranBarang.GET("/", storage.GetAllSebaranBarang)
		sebaranBarang.PUT("/:id", storage.UpdateSebaranBarang)
		sebaranBarang.DELETE("/:id", storage.DeleteSebaranBarang)
	}

	// User Routes
	user := r.Group("/user")
	{
		user.POST("/", storage.CreateUser)
		user.GET("/:id", storage.GetUserByID)
		user.PUT("/:id", storage.UpdateUser)
		user.DELETE("/:id", storage.DeleteUser)
	}
	
	return r
}
