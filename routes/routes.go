package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes mendaftarkan seluruh endpoint API dan handler global
func SetupRoutes(r *gin.Engine) {
	// Layanan Static File untuk gambar yang diunggah (/uploads/products/...)
	r.Static("/uploads", "./uploads")

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

	// Grup Utama API (/api)
	api := r.Group("/api")
	{
		AuthRoutes(api)
		ProductRoutes(api)
		RoleRoutes(api)
		UserRoutes(api)
	}
}
