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
	err := db.AutoMigrate( &types.Gudang{},
        &types.Inventaris{},
        &types.Divisi{},
        &types.User{},
        &types.SebaranBarang{},
        &types.Kategori{},
        &types.Depresiasi{},
		&types.History{},
		&types.JadwalDepresiasi{},
		&types.BarangStatus{},
		&types.BarangFoto{},
    )
	if err != nil {
		log.Fatal("Gagal melakukan migrasi:", err)
	}
	log.Println("Migrasi database selesai")
}
