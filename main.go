package main

import (
	"inventory/src/function/controller"
	"inventory/src/function/database"
	"inventory/src/function/routes"
	"github.com/gin-contrib/cors"

	"time"
	"github.com/gin-gonic/gin"
)



func main() {
	
	
	database.ConnectDB()
	database.MigrateDB()
	controller.InitiateScheduler()
	
	r := gin.Default()

	r.Static("/storage", "./storage")
	// Pasang middleware CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Pasang semua route dari folder routes
	routes.SetupRouter(r)

	// Jalankan server
	r.Run(":8080")
}
