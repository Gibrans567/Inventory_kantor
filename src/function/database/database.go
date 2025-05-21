package database

 import (
 	"database/sql"
 	"fmt"
 	"log"
 	"os"
 	"time"

 	_"github.com/go-sql-driver/mysql"
 	"gorm.io/driver/mysql"
 	"gorm.io/gorm"
 )

 var DB *gorm.DB

 func ConnectDB() {
 	dbUser := getEnv("DB_USER", "root")
 	dbPassword := getEnv("DB_PASSWORD", "")
 	dbHost := getEnv("DB_HOST", "mysql")
 	dbPort := getEnv("DB_PORT", "3306")
 	dbName := getEnv("DB_NAME", "inventory_kantor")

 	dsnRoot := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPassword, dbHost, dbPort)
 	sqlDB, err := sql.Open("mysql", dsnRoot)
 	if err != nil {
 		log.Fatalf("Gagal konek ke MySQL: %v", err)
 	}
 	defer sqlDB.Close()

 	err = sqlDB.Ping()
 	if err != nil {
 		log.Fatalf("Gagal ping MySQL: %v", err)
 	}

 	_, err = sqlDB.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
 	if err != nil {
 		log.Fatalf("Gagal membuat database: %v", err)
 	}
 	log.Println("Database berhasil dibuat atau sudah ada.")

 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
 		dbUser, dbPassword, dbHost, dbPort, dbName)

 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
 	if err != nil {
 		log.Fatalf("Gagal konek ke database: %v", err)
 	}

 	DB = db
 	log.Println("Koneksi ke database berhasil.")

 	for i := 0; i < 10; i++ {  
 		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
 		if err == nil {
 			DB = db
 			log.Println("Koneksi ke database berhasil.")
 			return
 		}
 		log.Printf("Gagal konek ke database, coba ulangi (%d/10): %v", i+1, err)
 		time.Sleep(3 * time.Second)
 	}

 	log.Fatalf("Gagal konek ke database setelah 10 percobaan: %v", err)
 }

 func GetDB() *gorm.DB {
 	if DB == nil {
 		log.Fatal("Database belum terkoneksi, pastikan ConnectDB() sudah dipanggil")
 	}
 	return DB
 }

 func getEnv(key, defaultVal string) string {
 	val := os.Getenv(key)
 	if val == "" {
 		return defaultVal
 	}
 	return val
 }
