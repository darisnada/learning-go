package middlewares

import (
	"net/http"

	"go-learning/config"
	"go-learning/models"

	"github.com/gin-gonic/gin"
)

// RequirePermission memeriksa apakah user yang login memiliki permission yang ditentukan
func RequirePermission(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Autentikasi diperlukan"})
			c.Abort()
			return
		}

		userIDStr, ok := userID.(string)
		if !ok || userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID tidak valid"})
			c.Abort()
			return
		}

		var user models.User
		if err := config.DB.Preload("Role.Permissions").Where("id = ?", userIDStr).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Pengguna tidak ditemukan"})
			c.Abort()
			return
		}

		if !user.HasPermission(permissionName) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":      "Akses ditolak: Anda tidak memiliki izin untuk tindakan ini",
				"permission": permissionName,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
