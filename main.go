package main

import (
	"inventory/src/function/controller"
	"inventory/src/function/database"
	"inventory/src/function/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	database.ConnectDB()
	database.MigrateDB()
	controller.InitiateScheduler()

	r := gin.Default()

	// Middleware CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	// Route backend API
	routes.SetupRouter(r)

	// Sajikan storage statis (misal untuk file upload)
	r.Static("/storage", "./storage")

	// Sajikan semua file static Angular dari folder browser
	r.StaticFS("/assets", http.Dir("./Inventaris/browser/assets")) // optional untuk folder assets Angular

	// Tangani route yang tidak cocok (SPA support)
	r.NoRoute(func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		// Cek apakah file statis ada
		filePath := filepath.Join("./Inventaris/browser", requestPath)

		if _, err := os.Stat(filePath); err == nil && !isDir(filePath) {
			c.File(filePath)
			return
		}

		// Fallback ke index.html untuk Angular routing
		c.File("./Inventaris/browser/index.html")
	})

	// Jalankan server
	r.Run(":8081")
}

// Helper untuk cek apakah path adalah direktori
func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
