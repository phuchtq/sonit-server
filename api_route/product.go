package apiroute

import (
	"sonit_server/handler"
	"sonit_server/utils/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeProductHandlerRoute(server *gin.Engine, port string) {
	// Context path
	var contextPath string = "products"

	// Define Product endpoints with admin required
	var norGroup = server.Group(contextPath)
	norGroup.GET("/customer-ui", handler.GetProductsCustomerUI)
	norGroup.GET("/price-interval", handler.GetProductsByPriceInterval)
	norGroup.GET("/:id", handler.GetProductById)

	var adminAuthGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorization)
	adminAuthGroup.GET("", handler.GetAllProducts)
	adminAuthGroup.GET("/category/:id", handler.GetProductsByCategory)
	adminAuthGroup.GET("/collection/:id", handler.GetProductsByCollection)
	adminAuthGroup.GET("/name/:name", handler.GetProductsByName)
	adminAuthGroup.GET("/description/:description", handler.GetProductsByDescription)
	adminAuthGroup.GET("/status/:status", handler.GetProductsByStatus)

	adminAuthGroup.POST("/create", handler.CreateProduct)
	adminAuthGroup.PUT("/update", handler.UpdateProduct)
	adminAuthGroup.DELETE("/remove", handler.RemoveProduct)
	adminAuthGroup.PATCH("/activate", handler.ActivateProduct)
}
