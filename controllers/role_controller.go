package controllers

import (
	"net/http"

	"go-learning/config"
	"go-learning/models"
	"go-learning/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RoleInput struct {
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description"`
	PermissionIDs []string `json:"permission_ids"` // Opsional: daftar UUID permission yang langsung dihubungkan
}

type UpdateRolePermissionsInput struct {
	PermissionIDs []string `json:"permission_ids" binding:"required"` // Daftar UUID permission baru
}

// GetRoles mengambil seluruh daftar role beserta permission masing-masing
func GetRoles(c *gin.Context) {
	var roles []models.Role
	if err := config.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data role: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil data role",
		"data":    roles,
	})
}

// GetRoleByID mengambil detail satu role berdasarkan ID UUID
func GetRoleByID(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid (harus berupa UUID)")
		return
	}

	var role models.Role
	if err := config.DB.Preload("Permissions").Where("id = ?", id).First(&role).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Role tidak ditemukan")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Detail role ditemukan",
		"data":    role,
	})
}

// CreateRole membuat role baru dan dapat langsung mengaitkan permissions
func CreateRole(c *gin.Context) {
	var input RoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Cek apakah nama role sudah ada
	var existing models.Role
	if err := config.DB.Where("name = ?", input.Name).First(&existing).Error; err == nil {
		utils.ErrorResponse(c, http.StatusConflict, "Nama role sudah digunakan")
		return
	}

	var perms []models.Permission
	if len(input.PermissionIDs) > 0 {
		config.DB.Where("id IN ?", input.PermissionIDs).Find(&perms)
	}

	role := models.Role{
		Name:        input.Name,
		Description: input.Description,
		Permissions: perms,
	}

	if err := config.DB.Create(&role).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal membuat role: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Role berhasil dibuat",
		"data":    role,
	})
}

// UpdateRole memperbarui informasi nama atau deskripsi role
func UpdateRole(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid (harus berupa UUID)")
		return
	}

	var role models.Role
	if err := config.DB.Where("id = ?", id).First(&role).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Role tidak ditemukan")
		return
	}

	var input RoleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Cegah penggantian nama role bawaan 'admin'
	if role.Name == "admin" && input.Name != "admin" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Nama role sistem 'admin' tidak boleh diubah")
		return
	}

	// Cek duplikasi nama jika nama diubah
	if input.Name != role.Name {
		var duplicate models.Role
		if err := config.DB.Where("name = ? AND id != ?", input.Name, id).First(&duplicate).Error; err == nil {
			utils.ErrorResponse(c, http.StatusConflict, "Nama role sudah digunakan oleh role lain")
			return
		}
	}

	role.Name = input.Name
	role.Description = input.Description

	if err := config.DB.Save(&role).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui role: "+err.Error())
		return
	}

	config.DB.Preload("Permissions").First(&role, "id = ?", role.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Role berhasil diperbarui",
		"data":    role,
	})
}

// UpdateRolePermissions memperbarui/mengganti seluruh permissions yang dimiliki suatu role
func UpdateRolePermissions(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid (harus berupa UUID)")
		return
	}

	var role models.Role
	if err := config.DB.Where("id = ?", id).First(&role).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Role tidak ditemukan")
		return
	}

	var input UpdateRolePermissionsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.ValidationErrorResponse(c, err)
		return
	}

	// Ambil daftar permissions berdasarkan UUID yang dikirim
	var perms []models.Permission
	if len(input.PermissionIDs) > 0 {
		config.DB.Where("id IN ?", input.PermissionIDs).Find(&perms)
	}

	// Perbarui relasi many-to-many role_permissions
	if err := config.DB.Model(&role).Association("Permissions").Replace(perms); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui permissions role: "+err.Error())
		return
	}

	// Ambil ulang role dengan permissions yang baru
	config.DB.Preload("Permissions").First(&role, "id = ?", role.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Permissions untuk role berhasil diperbarui",
		"data":    role,
	})
}

// DeleteRole menghapus role dari sistem
func DeleteRole(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Format ID tidak valid (harus berupa UUID)")
		return
	}

	var role models.Role
	if err := config.DB.Where("id = ?", id).First(&role).Error; err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Role tidak ditemukan")
		return
	}

	// Cegah penghapusan role sistem utama
	if role.Name == "admin" || role.Name == "user" {
		utils.ErrorResponse(c, http.StatusBadRequest, "Role sistem bawaan ('"+role.Name+"') tidak dapat dihapus")
		return
	}

	// Cek apakah ada pengguna yang masih menggunakan role ini
	var userCount int64
	config.DB.Model(&models.User{}).Where("role_id = ?", role.ID).Count(&userCount)
	if userCount > 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "Tidak dapat menghapus role yang masih digunakan oleh pengguna aktif")
		return
	}

	// Hapus relasi di tabel pivot terlebih dahulu
	config.DB.Model(&role).Association("Permissions").Clear()

	if err := config.DB.Delete(&role).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus role: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role berhasil dihapus",
	})
}

// GetPermissions mengambil seluruh daftar permission yang tersedia di sistem
func GetPermissions(c *gin.Context) {
	var perms []models.Permission
	if err := config.DB.Find(&perms).Error; err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil data permission: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Berhasil mengambil daftar permission",
		"data":    perms,
	})
}

