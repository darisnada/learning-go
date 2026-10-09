package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ErrorResponse mengembalikan respon error standar JSON
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"error": message,
	})
}

// ValidationErrorResponse memformat pesan error validasi agar mudah dibaca oleh client/frontend
func ValidationErrorResponse(c *gin.Context, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		out := make(map[string]string)
		for _, fe := range ve {
			field := fe.Field()
			switch fe.Tag() {
			case "required":
				out[field] = fmt.Sprintf("Field '%s' wajib diisi", field)
			case "email":
				out[field] = "Format email tidak valid"
			case "min":
				out[field] = fmt.Sprintf("Field '%s' minimal harus %s karakter", field, fe.Param())
			case "max":
				out[field] = fmt.Sprintf("Field '%s' maksimal %s karakter", field, fe.Param())
			case "gt":
				out[field] = fmt.Sprintf("Field '%s' harus lebih besar dari %s", field, fe.Param())
			case "gte":
				out[field] = fmt.Sprintf("Field '%s' minimal harus %s", field, fe.Param())
			default:
				out[field] = fmt.Sprintf("Field '%s' tidak memenuhi validasi (%s)", field, fe.Tag())
			}
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validasi input gagal",
			"details": out,
		})
		return
	}

	// Jika format payload tidak valid atau bukan validation error
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Payload request tidak valid: " + err.Error(),
	})
}

