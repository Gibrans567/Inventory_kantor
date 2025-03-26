package database

import (
	"log"
	"inventory/src/types"
)

func MigrateDB() {
	db := GetDB() // Ambil koneksi database
	if db == nil {
		log.Fatal("Database belum terkoneksi, pastikan ConnectDB() sudah dipanggil")
	}

	// Lakukan migrasi
	err := db.AutoMigrate(&types.Inventory{})
	if err != nil {
		log.Fatal("Gagal melakukan migrasi:", err)
	}
	log.Println("Migrasi database selesai")
}
