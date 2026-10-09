package controllers

import (
	"net/http"

	"go-learning/config"
	"go-learning/models"
	"go-learning/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Stock       int     `json:"stock" binding:"gte=0"`
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail produk ditemukan",
		"data":    product,
	})
}

// CreateProduct menambahkan produk baru (memerlukan token login)
func CreateProduct(c *gin.Context) {
	var input ProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
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

	product := models.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Stock:       input.Stock,
		UserID:      userIDStr,
	}

	if err := config.DB.Create(&product).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan data produk: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Produk berhasil ditambahkan",
		"data":    product,
	})
}

// UpdateProduct memperbarui data produk berdasarkan ID UUID
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
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	product.Name = input.Name
	product.Description = input.Description
	product.Price = input.Price
	product.Stock = input.Stock

	if err := config.DB.Save(&product).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui data produk: "+err.Error())
		return
	}

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

	if err := config.DB.Delete(&product).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus produk: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Produk berhasil dihapus",
	})
}
