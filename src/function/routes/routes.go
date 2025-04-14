package routes

import (
	"inventory/src/function/controller"
	"github.com/gin-gonic/gin"
	"inventory/src/function/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Depresiasi Routes
	depresiasi := r.Group("/depresiasi")
	{
		depresiasi.POST("", controller.CreateDepresiasi)
		depresiasi.GET("", controller.GetAllDepresiasi)
		depresiasi.PUT("/:id", controller.UpdateDepresiasi)
		depresiasi.DELETE("/:id", controller.DeleteDepresiasi)
	}

	// Divisi Routes
	divisi := r.Group("/divisi")
	{
		divisi.POST("", controller.CreateDivisi)
		divisi.GET("", controller.GetAllDivisi)
		divisi.PUT("/:id", controller.UpdateDivisi)
		divisi.DELETE("/:id", controller.DeleteDivisi)
	}

	// Gudang Routes
	gudang := r.Group("/gudang")
	{
		gudang.POST("", controller.CreateGudang)
		gudang.GET("", controller.GetAllGudang)
		gudang.POST("/update/:id", controller.UpdateGudang)
		gudang.DELETE("/:id", controller.DeleteGudang)
	}

	// Inventaris Routes
	inventaris := r.Group("/inventaris")
	{
		inventaris.POST("", controller.CreateInventaris)
		inventaris.GET("", controller.GetAllInventaris)
		inventaris.GET("/:nama_divisi", controller.GetInventarisByDivisiName)
		inventaris.PUT("/:id", controller.UpdateInventaris)
		inventaris.DELETE("/:id", controller.DeleteInventaris)
	}

	// Kategori Routes
	kategori := r.Group("/kategori")
	{
		kategori.POST("", controller.CreateKategori)
		kategori.GET("", controller.GetAllKategori)
		kategori.PUT("/:id", controller.UpdateKategori)
		kategori.DELETE("/:id", controller.DeleteKategori)
	}

	// SebaranBarang Routes
	sebaranBarang := r.Group("/sebaranBarang")
	{
		sebaranBarang.POST("", controller.CreateSebaranBarang)
		sebaranBarang.GET("", controller.GetAllSebaranBarang)
		sebaranBarang.PUT("/:id", controller.UpdateSebaranBarang)
		sebaranBarang.DELETE("/:id", controller.DeleteSebaranBarang)
	}

	// User Routes
	user := r.Group("/user")
	{
		user.POST("", controller.CreateUser)
		user.GET("", controller.GetUserByID)
		user.GET("/:id", controller.GetUserByID)
		user.PUT("/:id", controller.UpdateUser)
		user.DELETE("/:id", controller.DeleteUser)
	}
	
	historyRoutes := r.Group("/histories")
	{
		historyRoutes.POST("", controller.CreateHistory)
		historyRoutes.GET("", controller.GetAllHistories)
		historyRoutes.GET("/:id", controller.GetHistoryByID)
		historyRoutes.POST("/:id", controller.UpdateHistory)
		historyRoutes.DELETE("/:id", controller.DeleteHistory)
	}

	schedulerRoutes := r.Group("/scheduler")
	{
		schedulerRoutes.POST("", controller.ApplyDepresiasi)
		schedulerRoutes.GET("", controller.GetAllJadwal)
	}

	DeleteRoutes := r.Group("/delete")
	{
		DeleteRoutes.POST("", controller.ApplyDepresiasi)
		DeleteRoutes.DELETE("", controller.DeleteAllByTimeframe)
	}

	uploadRoutes := r.Group("/upload")
	{
		uploadRoutes.POST("", controller.UploadGambar)
	}


	// Auth group
	auth := r.Group("/auth")
	{
		auth.POST("/login", middleware.LoginHandler)
	}

	// Protected route
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/pegawai", middleware.RequireRole("pegawai"),  middleware.LogoutHandler)
	}

	return r
}
