package controllers

import (
	"net/http"

	"go-learning/config"
	"go-learning/models"
	"go-learning/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	RoleID   string `json:"role_id" binding:"required"`
}

type UpdateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"omitempty,min=6"` // Opsional: hanya jika ingin reset kata sandi
	RoleID   string `json:"role_id"`                            // Opsional
}

type UpdateUserRoleInput struct {
	RoleID string `json:"role_id" binding:"required"`
}

// GetUsers mengambil seluruh daftar pengguna dengan filter pencarian nama/email
func GetUsers(c *gin.Context) {
	var users []models.User

	query := config.DB.Preload("Role.Permissions")
	search := c.Query("search")
	if search != "" {
		query = query.Where("name LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&users).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data pengguna: "+err.Error())
		return
	}

	// Bersihkan password hash dari respon
	for i := range users {
		users[i].Password = ""
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data pengguna",
		"data":    users,
	})
}

// GetUserByID mengambil detail satu pengguna berdasarkan ID UUID
func GetUserByID(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid (harus berupa UUID)")
		return
	}

	var user models.User
	if err := config.DB.Preload("Role.Permissions").Where("id = ?", id).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Pengguna tidak ditemukan")
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail pengguna ditemukan",
		"data":    user,
	})
}

// CreateUser membuat pengguna baru oleh admin beserta penentuan role-nya
func CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Cek apakah email sudah terdaftar
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Email sudah terdaftar")
		return
	}

	// Cek apakah RoleID valid dan ada di database
	if _, err := uuid.Parse(input.RoleID); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format Role ID tidak valid (harus berupa UUID)")
		return
	}

	var role models.Role
	if err := config.DB.Where("id = ?", input.RoleID).First(&role).Error; err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Role dengan ID tersebut tidak ditemukan")
		return
	}

	user := models.User{
		Name:   input.Name,
		Email:  input.Email,
		RoleID: role.ID,
	}

	if err := user.HashPassword(input.Password); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengenkripsi kata sandi: "+err.Error())
		return
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat pengguna: "+err.Error())
		return
	}

	config.DB.Preload("Role.Permissions").First(&user, "id = ?", user.ID)
	user.Password = ""

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pengguna baru berhasil dibuat",
		"data":    user,
	})
}

// UpdateUser memperbarui data profil pengguna (nama, email, dan opsional password/role)
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid (harus berupa UUID)")
		return
	}

	var user models.User
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Pengguna tidak ditemukan")
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Cek jika email diubah dan sudah dipakai akun lain
	if input.Email != user.Email {
		var duplicate models.User
		if err := config.DB.Where("email = ? AND id != ?", input.Email, id).First(&duplicate).Error; err == nil {
			utils.ErrorResponse(c, http.StatusConflict, "Email sudah digunakan oleh pengguna lain")
			return
		}
	}

	user.Name = input.Name
	user.Email = input.Email

	// Jika ada pergantian password
	if input.Password != "" {
		if err := user.HashPassword(input.Password); err != nil {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengenkripsi kata sandi baru: "+err.Error())
			return
		}
	}

	// Jika ada pergantian role
	if input.RoleID != "" {
		var role models.Role
		if err := config.DB.Where("id = ?", input.RoleID).First(&role).Error; err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Role dengan ID tersebut tidak ditemukan")
			return
		}
		user.RoleID = role.ID
	}

	if err := config.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui pengguna: "+err.Error())
		return
	}

	config.DB.Preload("Role.Permissions").First(&user, "id = ?", user.ID)
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "Data pengguna berhasil diperbarui",
		"data":    user,
	})
}

// UpdateUserRole khusus untuk mengubah / menetapkan role seorang user
func UpdateUserRole(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid (harus berupa UUID)")
		return
	}

	var user models.User
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Pengguna tidak ditemukan")
		return
	}

	var input UpdateUserRoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Cek apakah RoleID valid di database
	var role models.Role
	if err := config.DB.Where("id = ?", input.RoleID).First(&role).Error; err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Role dengan ID tersebut tidak ditemukan")
		return
	}

	user.RoleID = role.ID
	if err := config.DB.Save(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui role pengguna: "+err.Error())
		return
	}

	config.DB.Preload("Role.Permissions").First(&user, "id = ?", user.ID)
	user.Password = ""

	c.JSON(http.StatusOK, gin.H{
		"message": "Role pengguna berhasil diperbarui",
		"data":    user,
	})
}

// DeleteUser menghapus pengguna dari sistem
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid (harus berupa UUID)")
		return
	}

	// Proteksi: cegah user menghapus dirinya sendiri yang sedang login
	currentUserID, exists := c.Get("user_id")
	if exists && currentUserID.(string) == id {
		utils.ErrorResponse(c, http.StatusBadRequest, "Anda tidak dapat menghapus akun Anda sendiri yang sedang aktif")
		return
	}

	var user models.User
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Pengguna tidak ditemukan")
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus pengguna: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pengguna berhasil dihapus",
	})
}

