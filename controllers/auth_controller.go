package controllers

import (
	"net/http"

	"go-learning/config"
	"go-learning/models"
	"go-learning/utils"

	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   string `json:"role_id"` // Opsional, default role jika kosong adalah "user"
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register membuat akun user baru
func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah email sudah terdaftar
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email sudah terdaftar"})
		return
	}

	// Tentukan RoleID (default ke role 'user' jika tidak disediakan)
	roleID := input.RoleID
	if roleID == "" {
		var userRole models.Role
		if err := config.DB.Where("name = ?", "user").First(&userRole).Error; err == nil {
			roleID = userRole.ID
		}
	}

	user := models.User{
		Name:   input.Name,
		Email:  input.Email,
		RoleID: roleID,
	}

	if err := user.HashPassword(input.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengenkripsi kata sandi"})
		return
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pengguna: " + err.Error()})
		return
	}

	// Preload role untuk response
	config.DB.Preload("Role.Permissions").First(&user, "id = ?", user.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registrasi berhasil",
		"data": gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		},
	})
}

// Login melakukan otentikasi user dan mengembalikan JWT token
func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Preload("Role.Permissions").Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau kata sandi salah"})
		return
	}

	if err := user.CheckPassword(input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau kata sandi salah"})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghasilkan token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

// GetProfile mengambil data profil user yang sedang login beserta role & permissions
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak ditemukan pada sesi"})
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Format User ID tidak valid"})
		return
	}

	var user models.User
	if err := config.DB.Preload("Role.Permissions").Where("id = ?", userIDStr).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"created_at": user.CreatedAt,
		},
	})
}
