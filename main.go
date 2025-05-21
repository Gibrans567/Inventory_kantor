package main

import (
	"inventory/src/function/controller"
	"inventory/src/function/database"
	"inventory/src/function/routes"
	"github.com/gin-contrib/cors"
	

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
		AllowOrigins:     []string{"*"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	// Pasang semua route dari folder routes
	routes.SetupRouter(r)

	// Jalankan server
	r.Run(":8081")
}
