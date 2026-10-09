package controllers

import (
	"errors"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go-learning/config"
	"go-learning/models"
	"go-learning/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductInput struct {
	Name        string                `form:"name" json:"name" binding:"required"`
	Description string                `form:"description" json:"description"`
	Price       float64               `form:"price" json:"price" binding:"required,gt=0"`
	Stock       int                   `form:"stock" json:"stock" binding:"gte=0"`
	Image       *multipart.FileHeader `form:"image" json:"-"`
	ImageURL    string                `form:"imageUrl" json:"image"` // Opsional: URL/path string jika dikirim via JSON
}

// handleImageUpload memproses dan memvalidasi file gambar dari request multipart/form-data
func handleImageUpload(c *gin.Context, fileHeader *multipart.FileHeader) (string, error) {
	file := fileHeader
	if file == nil {
		var err error
		file, err = c.FormFile("image")
		if err != nil {
			// Tidak ada file gambar yang diunggah
			return "", nil
		}
	}

	// Validasi ukuran gambar (maksimal 5 MB)
	if file.Size > 5*1024*1024 {
		return "", errors.New("ukuran gambar melebihi batas maksimal 5MB")
	}

	// Validasi ekstensi gambar
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".gif":  true,
	}

	if !allowedExtensions[ext] {
		return "", errors.New("format gambar tidak didukung (hanya .jpg, .jpeg, .png, .webp, .gif)")
	}

	// Buat direktori upload jika belum ada
	uploadDir := "./uploads/products"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", errors.New("gagal membuat direktori upload gambar")
	}

	// Generate nama file unik dengan UUID v4
	newFileName := uuid.NewString() + ext
	dst := filepath.Join(uploadDir, newFileName)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", errors.New("gagal menyimpan file gambar ke server")
	}

	// Return path URL yang dapat diakses publik
	return "/uploads/products/" + newFileName, nil
}

// formatImageURL mengubah path relatif gambar menjadi full URL publik (misal: http://localhost:8080/uploads/products/xxx.png)
func formatImageURL(c *gin.Context, imagePath string) string {
	if imagePath == "" {
		return ""
	}
	if strings.HasPrefix(imagePath, "http://") || strings.HasPrefix(imagePath, "https://") {
		return imagePath
	}

	appURL := os.Getenv("APP_URL")
	if appURL != "" {
		return strings.TrimRight(appURL, "/") + imagePath
	}

	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	host := c.Request.Host
	if host == "" {
		host = "localhost:8080"
	}
	return scheme + "://" + host + imagePath
}

// formatProduct mengubah field Image pada Product menjadi full URL
func formatProduct(c *gin.Context, product *models.Product) {
	if product.Image != "" {
		product.Image = formatImageURL(c, product.Image)
	}
}

// GetProducts menampilkan seluruh daftar produk
func GetProducts(c *gin.Context) {
	var products []models.Product

	// Opsi pencarian jika ada query parameter "search"
	query := config.DB.Preload("User")
	search := c.Query("search")
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	if err := query.Find(&products).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data produk: "+err.Error())
		return
	}

	for i := range products {
		formatProduct(c, &products[i])
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data produk",
		"data":    products,
	})
}

// GetProductByID menampilkan detail satu produk berdasarkan ID UUID
func GetProductByID(c *gin.Context) {
	id := c.Param("id")

	// Validasi format UUID
	if _, err := uuid.Parse(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid (harus berupa UUID)")
		return
	}

	var product models.Product
	if err := config.DB.Preload("User").Where("id = ?", id).First(&product).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Produk dengan ID tersebut tidak ditemukan")
		return
	}

	formatProduct(c, &product)

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail produk ditemukan",
		"data":    product,
	})
}

// CreateProduct menambahkan produk baru (mendukung upload file gambar via multipart/form-data atau JSON)
func CreateProduct(c *gin.Context) {
	var input ProductInput
	if err := c.ShouldBind(&input); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Autentikasi diperlukan")
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		utils.ErrorResponse(c, http.StatusUnauthorized, "Sesi User ID tidak valid")
		return
	}

	// Tangani upload gambar
	imagePath, err := handleImageUpload(c, input.Image)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Jika tidak upload file tapi mengirim URL/path di JSON
	if imagePath == "" && input.ImageURL != "" {
		imagePath = input.ImageURL
	}

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		Image:       imagePath,
		UserID:      userIDStr,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan data produk: "+err.Error())
		return
	}

	formatProduct(c, &product)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Produk berhasil ditambahkan",
		"data":    product,
	})
}

// UpdateProduct memperbarui data produk berdasarkan ID UUID (mendukung ganti gambar)
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	// Validasi format UUID
	if _, err := uuid.Parse(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid (harus berupa UUID)")
		return
	}

	var product models.Product
	if err := config.DB.Where("id = ?", id).First(&product).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Produk dengan ID tersebut tidak ditemukan")
		return
	}

	var input ProductInput
	if err := c.ShouldBind(&input); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Tangani upload gambar baru jika ada
	newImagePath, err := handleImageUpload(c, input.Image)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Stock = input.Stock

	// Jika upload gambar baru, hapus gambar lama dari disk jika ada
	if newImagePath != "" {
		if product.Image != "" && strings.HasPrefix(product.Image, "/uploads/products/") {
			_ = os.Remove("." + product.Image)
		}
		product.Image = newImagePath
	} else if input.ImageURL != "" {
		product.Image = input.ImageURL
	}

	if err := config.DB.Save(&product).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui data produk: "+err.Error())
		return
	}

	formatProduct(c, &product)

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk berhasil diperbarui",
		"data":    product,
	})
}

// DeleteProduct menghapus produk berdasarkan ID UUID
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	// Validasi format UUID
	if _, err := uuid.Parse(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid (harus berupa UUID)")
		return
	}

	var product models.Product
	if err := config.DB.Where("id = ?", id).First(&product).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Produk dengan ID tersebut tidak ditemukan")
		return
	}

	// Hapus file gambar jika ada
	if product.Image != "" && strings.HasPrefix(product.Image, "/uploads/products/") {
		_ = os.Remove("." + product.Image)
	}

	if err := config.DB.Delete(&product).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus produk: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk berhasil dihapus",
	})
}
