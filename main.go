package main

import (
	"inventory/src/function/controller"
	"inventory/src/function/database"
	"inventory/src/function/routes"
)



func main() {
	
	
	database.ConnectDB()
	database.MigrateDB()
	controller.InitiateScheduler()
	
	r := routes.SetupRouter()

	// Jalankan server
	r.Run(":8080")
}
