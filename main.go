package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"go-learning/config"
	"go-learning/controllers"
	"go-learning/middlewares"
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

	// 5. Inisialisasi router Gin
	r := gin.Default()

	// Penanganan Error Global 404 Route Not Found dalam format JSON
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Endpoint URL tidak ditemukan",
		})
	})

	// Penanganan Error Global 405 Method Not Allowed dalam format JSON
	r.NoMethod(func(c *gin.Context) {
		c.JSON(http.StatusMethodNotAllowed, gin.H{
			"error": "Metode HTTP tidak diizinkan untuk endpoint ini",
		})
	})

	// Root / Health Check
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "API Go + Gin + GORM + MySQL (RBAC Active) berjalan normal!",
			"status":  "healthy",
		})
	})

	api := r.Group("/api")
	{
		// Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			// Endpoint profile memerlukan token login
			auth.GET("/me", middlewares.AuthMiddleware(), controllers.GetProfile)
		}

		// Product Routes (Dilindungi oleh Auth & Permission RBAC)
		products := api.Group("/products")
		products.Use(middlewares.AuthMiddleware())
		{
			// Izin: read-product
			products.GET("", middlewares.RequirePermission("read-product"), controllers.GetProducts)
			products.GET("/:id", middlewares.RequirePermission("read-product"), controllers.GetProductByID)

			// Izin: create-product
			products.POST("", middlewares.RequirePermission("create-product"), controllers.CreateProduct)

			// Izin: update-product
			products.PUT("/:id", middlewares.RequirePermission("update-product"), controllers.UpdateProduct)

			// Izin: delete-product
			products.DELETE("/:id", middlewares.RequirePermission("delete-product"), controllers.DeleteProduct)
		}
	}

	// 6. Jalankan server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server berjalan di http://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
