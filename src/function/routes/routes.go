package routes

import (
	"github.com/gin-gonic/gin"
	"inventory/src/function/controller"
	"inventory/src/function/middleware"
)

func SetupRouter(r *gin.Engine){

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
		divisi.DELETE("/:nama_divisi", controller.DeleteDivisi)
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
		inventaris.GET("/barang/:id", controller.GetInventarisById)
		inventaris.PUT("/:id", controller.UpdateInventaris)
		inventaris.DELETE("/:id", controller.DeleteInventaris)
	}

	// Kategori Routes
	kategori := r.Group("/kategori")
	{
		kategori.POST("", controller.CreateKategori)
		kategori.GET("", controller.GetAllKategori)
		kategori.PUT("/:id", controller.UpdateKategori)
		kategori.DELETE("/:nama_kategori", controller.DeleteKategori)
	}

	// SebaranBarang Routes
	sebaranBarang := r.Group("/sebaranBarang")
	{
		sebaranBarang.POST("", controller.CreateSebaranBarang)
		sebaranBarang.GET("", controller.GetAllSebaranBarang)
		sebaranBarang.GET("/sebaran/:id", controller.GetSebaranBarangByID)
		sebaranBarang.GET("/:id", controller.GetSebaranBarangByIDBarang)
		sebaranBarang.PUT("/Edit/:id", controller.UpdateSebaranBarang)
		sebaranBarang.POST("/pindah", controller.MoveSebaranBarang)
		sebaranBarang.DELETE("/:id", controller.DeleteSebaranBarang)
	}

	// User Routes
	user := r.Group("/user")
	{
		user.POST("", controller.CreateUser)
		user.GET("", controller.GetAllUsers)
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

	// RegisterBarangStatusRoutes registers all the barang status routes
	barangStatusRoutes := r.Group("/barangStatus")
	{
	barangStatusRoutes.POST("", controller.CreateBarangStatus)
	barangStatusRoutes.GET("", controller.GetAllBarangStatus)
	barangStatusRoutes.GET("/:id", controller.GetBarangStatus)
	barangStatusRoutes.GET("/barang/:id", controller.GetBarangStatusByBarang)
	barangStatusRoutes.GET("/sebaran/:id", controller.GetBarangStatusBySebaran)
	barangStatusRoutes.PUT("/:id", controller.UpdateBarangStatus)
	barangStatusRoutes.DELETE("/:id", controller.DeleteBarangStatus)
	}

	MultiUploadGambarRoutes := r.Group("/MultiUploadGambar")
	{
	MultiUploadGambarRoutes.POST("", controller.UploadGambarMulti)
	MultiUploadGambarRoutes.GET("", controller.GetAllBarangFoto)
	MultiUploadGambarRoutes.GET("/:id", controller.GetBarangFoto)
	MultiUploadGambarRoutes.GET("/barang/:id", controller.GetBarangFotoByBarang)
	MultiUploadGambarRoutes.POST("/:id", controller.UpdateBarangFoto)
	MultiUploadGambarRoutes.DELETE("/:id", controller.DeleteBarangFoto)
	MultiUploadGambarRoutes.DELETE("/barang/:id", controller.DeleteAllBarangFotoByBarang)
	}

	pengajuan := r.Group("/pengajuan")
	{ 
	pengajuan.POST("", controller.CreatePengajuan)
	pengajuan.GET("", controller.GetAllPengajuan)
	pengajuan.GET("/:id", controller.GetPengajuanByID)
	pengajuan.PUT("/:id", controller.UpdatePengajuan)
	pengajuan.DELETE("/:id", controller.DeletePengajuan)
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
		auth.POST("/register", middleware.RegisterHandler)
	}

	// Protected route
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.POST("/pegawai", middleware.RequireRole("pegawai"), middleware.LogoutHandler)
	}

}
