package main

import (
	"flag"
	"log"
	"os"

	"go-learning/config"
	"go-learning/routes"
	"go-learning/seeders"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Muat variabel lingkungan dari file .env jika ada
	if err := godotenv.Load(); err != nil {
		log.Println("Info: File .env tidak ditemukan, menggunakan environment variable default")
	}

	// 2. Definisi CLI Flags untuk Migration & Seeder
	migrateOnly := flag.Bool("migrate", false, "Hanya jalankan auto-migrate tabel lalu keluar")
	seedOnly := flag.Bool("seed", false, "Jalankan seeder data dummy lalu keluar")
	freshOnly := flag.Bool("fresh", false, "Hapus tabel lama dan migrasi ulang dari awal (fresh)")
	freshSeed := flag.Bool("fresh-seed", false, "Hapus tabel, migrasi ulang, dan isi data seeder")
	flag.Parse()

	// 3. Hubungkan ke database MySQL
	config.ConnectDatabase()

	// 4. Tangani perintah CLI jika ada flag yang dipanggil
	if *freshOnly {
		config.Fresh()
		log.Println("Database fresh migration selesai!")
		return
	}

	if *freshSeed {
		config.Fresh()
		seeders.SeedAll(config.DB)
		log.Println("Database fresh migration dan seed selesai!")
		return
	}

	if *migrateOnly {
		config.Migrate()
		log.Println("Migrasi selesai!")
		return
	}

	if *seedOnly {
		config.Migrate()
		seeders.SeedAll(config.DB)
		log.Println("Seeding selesai!")
		return
	}

	// Default: selalu jalankan migrasi tabel sebelum menyalakan server
	config.Migrate()

	// 5. Inisialisasi router Gin dan daftarkan routing secara modular
	r := gin.Default()
	routes.SetupRoutes(r)

	// 6. Jalankan server HTTP
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server berjalan di http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
