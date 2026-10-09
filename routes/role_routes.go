package routes

import (
	"go-learning/controllers"
	"go-learning/middlewares"

	"github.com/gin-gonic/gin"
)

// RoleRoutes mendaftarkan routing untuk fitur manajemen role & permissions
func RoleRoutes(api *gin.RouterGroup) {
	roles := api.Group("/roles")
	roles.Use(middlewares.AuthMiddleware())
	{
		// Izin: read-role
		roles.GET("", middlewares.RequirePermission("read-role"), controllers.GetRoles)
		roles.GET("/:id", middlewares.RequirePermission("read-role"), controllers.GetRoleByID)

		// Izin: create-role
		roles.POST("", middlewares.RequirePermission("create-role"), controllers.CreateRole)

		// Izin: update-role
		roles.PUT("/:id", middlewares.RequirePermission("update-role"), controllers.UpdateRole)

		// Izin: update-role (khusus update daftar permissions)
		roles.PUT("/:id/permissions", middlewares.RequirePermission("update-role"), controllers.UpdateRolePermissions)

		// Izin: delete-role
		roles.DELETE("/:id", middlewares.RequirePermission("delete-role"), controllers.DeleteRole)
	}

	// Permissions List Route (Daftar semua permissions yang ada di sistem)
	api.GET("/permissions", middlewares.AuthMiddleware(), middlewares.RequirePermission("read-role"), controllers.GetPermissions)
}

