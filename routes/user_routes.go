package routes

import (
	"go-learning/controllers"
	"go-learning/middlewares"

	"github.com/gin-gonic/gin"
)

// UserRoutes mendaftarkan routing untuk fitur manajemen pengguna
func UserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	users.Use(middlewares.AuthMiddleware())
	{
		// Izin: read-user
		users.GET("", middlewares.RequirePermission("read-user"), controllers.GetUsers)
		users.GET("/:id", middlewares.RequirePermission("read-user"), controllers.GetUserByID)

		// Izin: create-user
		users.POST("", middlewares.RequirePermission("create-user"), controllers.CreateUser)

		// Izin: update-user
		users.PUT("/:id", middlewares.RequirePermission("update-user"), controllers.UpdateUser)

		// Izin: update-user (khusus update role)
		users.PUT("/:id/role", middlewares.RequirePermission("update-user"), controllers.UpdateUserRole)

		// Izin: delete-user
		users.DELETE("/:id", middlewares.RequirePermission("delete-user"), controllers.DeleteUser)
	}
}

