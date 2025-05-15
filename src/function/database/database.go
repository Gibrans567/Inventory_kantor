package database

import (
	"database/sql"
	"fmt"
	"os"
	"log"

	_"github.com/go-sql-driver/mysql" // Import driver MySQL
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Konfigurasi koneksi database
const (
	DBUser     = "root"
	DBPassword = ""
	DBPort     = "3306"
	DBName     = "inventory_kantor"
)

var DBHost = os.Getenv("DB_HOST")
// ConnectDB membuat database jika belum ada, lalu menghubungkan GORM ke MySQL.
func ConnectDB() {
	// 1. Buat koneksi awal tanpa memilih database
	dsnRoot := fmt.Sprintf("%s:%s@tcp(%s:%s)/", DBUser, DBPassword, DBHost, DBPort)
	sqlDB, err := sql.Open("mysql", dsnRoot)
	if err != nil {
		log.Fatalf("Gagal konek ke MySQL: %v", err)
	}
	defer sqlDB.Close()

	// 2. Buat database jika belum ada
	_, err = sqlDB.Exec("CREATE DATABASE IF NOT EXISTS " + DBName)
	if err != nil {
		log.Fatalf("Gagal membuat database: %v", err)
	}
	log.Println("Database berhasil dibuat atau sudah ada.")

	// 3. Buat koneksi ke database yang sudah dibuat
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		DBUser, DBPassword, DBHost, DBPort, DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal konek ke database: %v", err)
	}

	// 4. Simpan koneksi ke variabel global
	DB = db
	log.Println("Koneksi ke database berhasil.")
}

// GetDB mengembalikan instance database
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database belum terkoneksi, pastikan ConnectDB() sudah dipanggil")
	}
	return DB
}

func CheckMySQLVersion() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		DBUser, DBPassword, DBHost, DBPort, DBName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Gagal konek ke MySQL untuk cek versi: %v", err)
	}
	defer db.Close()

	var version string
	err = db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		log.Fatalf("Gagal mengambil versi MySQL: %v", err)
	}

	fmt.Println("Versi MySQL:", version)
}
