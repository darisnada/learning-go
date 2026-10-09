package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"go-learning/models"

	_ "github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase initializes MySQL connection
func ConnectDatabase() {
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "go_restapi"
	}

	// 1. Pastikan database dibuat otomatis jika belum ada
	serverDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort)
	initDB, err := sql.Open("mysql", serverDSN)
	if err == nil {
		_, err = initDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;", dbName))
		if err != nil {
			log.Printf("Peringatan: Gagal membuat database otomatis: %v", err)
		}
		initDB.Close()
	}

	// 2. Hubungkan ke database dengan GORM
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName)
	database, err := gorm.Open(gormmysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database MySQL: %v", err)
	}

	DB = database
	log.Println("Database MySQL berhasil terhubung!")
}

// Migrate menjalankan migrasi schema tabel
func Migrate() {
	err := DB.AutoMigrate(
		&models.Permission{},
		&models.Role{},
		&models.User{},
		&models.Product{},
	)
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi tabel: %v", err)
	}
	log.Println("Migrasi tabel (permissions, roles, users, products) berhasil!")
}

// Fresh menghapus tabel dan melakukan migrasi ulang dari awal (reset)
func Fresh() {
	log.Println("Mereset tabel lama...")
	DB.Migrator().DropTable(
		"role_permissions",
		&models.Product{},
		&models.User{},
		&models.Role{},
		&models.Permission{},
	)
	Migrate()
}
