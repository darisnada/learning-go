package routes

import (
	"go-learning/controllers"
	"go-learning/middlewares"

	"github.com/gin-gonic/gin"
)

// ProductRoutes mendaftarkan routing untuk fitur manajemen produk
func ProductRoutes(api *gin.RouterGroup) {
	products := api.Group("/products")
	products.Use(middlewares.AuthMiddleware())
	{
		// Izin: read-product
		products.GET("", middlewares.RequirePermission("read-product"), controllers.GetProducts)
		products.GET("/:id", middlewares.RequirePermission("read-product"), controllers.GetProductByID)

		// Izin: create-product
		products.POST("", middlewares.RequirePermission("create-product"), controllers.CreateProduct)

		// Izin: update-product
		products.PUT("/:id", middlewares.RequirePermission("update-product"), controllers.UpdateProduct)

		// Izin: delete-product
		products.DELETE("/:id", middlewares.RequirePermission("delete-product"), controllers.DeleteProduct)
	}
}

