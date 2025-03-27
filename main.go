package main

import (
	"inventory/src/function/database"
	"inventory/src/function/routes"
)

func main() {
	// Koneksi ke database
	
	database.ConnectDB()
	database.MigrateDB()

	// // Setup router
	r := routes.SetupRouter()

	// Jalankan server
	r.Run(":8080")
}
