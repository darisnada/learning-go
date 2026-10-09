package routes

import (
	"go-learning/controllers"
	"go-learning/middlewares"

	"github.com/gin-gonic/gin"
)

// AuthRoutes mendaftarkan routing untuk fitur autentikasi
func AuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.GET("/me", middlewares.AuthMiddleware(), controllers.GetProfile)
	}
}

