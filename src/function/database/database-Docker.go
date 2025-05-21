package database

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"
// 	"time"

// 	_ "github.com/go-sql-driver/mysql" // driver MySQL
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// var DB *gorm.DB

// // ConnectDB membuat database jika belum ada, lalu menghubungkan GORM ke MySQL.
// func ConnectDB() {
// 	// Ambil konfigurasi dari environment variable dengan default fallback
// 	dbUser := getEnv("DB_USER", "root")
// 	dbPassword := getEnv("DB_PASSWORD", "")
// 	dbHost := getEnv("DB_HOST", "mysql")
// 	dbPort := getEnv("DB_PORT", "3306")
// 	dbName := getEnv("DB_NAME", "inventory_kantor")

// 	// 1. Koneksi awal tanpa memilih database untuk create DB jika belum ada
// 	dsnRoot := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPassword, dbHost, dbPort)
// 	sqlDB, err := sql.Open("mysql", dsnRoot)
// 	if err != nil {
// 		log.Fatalf("Gagal konek ke MySQL: %v", err)
// 	}
// 	defer sqlDB.Close()

// 	// Pastikan koneksi bisa dilakukan
// 	err = sqlDB.Ping()
// 	if err != nil {
// 		log.Fatalf("Gagal ping MySQL: %v", err)
// 	}

// 	// 2. Buat database jika belum ada
// 	_, err = sqlDB.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
// 	if err != nil {
// 		log.Fatalf("Gagal membuat database: %v", err)
// 	}
// 	log.Println("Database berhasil dibuat atau sudah ada.")

// 	// 3. Koneksi ke database yang sudah dibuat
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		dbUser, dbPassword, dbHost, dbPort, dbName)

// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Gagal konek ke database: %v", err)
// 	}

// 	// 4. Simpan koneksi ke variabel global
// 	DB = db
// 	log.Println("Koneksi ke database berhasil.")

// 	for i := 0; i < 10; i++ { // retry 10x
// 		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 		if err == nil {
// 			DB = db
// 			log.Println("Koneksi ke database berhasil.")
// 			return
// 		}
// 		log.Printf("Gagal konek ke database, coba ulangi (%d/10): %v", i+1, err)
// 		time.Sleep(3 * time.Second)
// 	}

// 	log.Fatalf("Gagal konek ke database setelah 10 percobaan: %v", err)
// }

// // GetDB mengembalikan instance database
// func GetDB() *gorm.DB {
// 	if DB == nil {
// 		log.Fatal("Database belum terkoneksi, pastikan ConnectDB() sudah dipanggil")
// 	}
// 	return DB
// }

// // Utility ambil env dengan default
// func getEnv(key, defaultVal string) string {
// 	val := os.Getenv(key)
// 	if val == "" {
// 		return defaultVal
// 	}
// 	return val
// }
