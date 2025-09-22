package apiroute

import (
	"sonit_server/handler"
	"sonit_server/utils/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeProductInventoryTransactionHandlerRoute(server *gin.Engine, port string) {
	var adminAuthGroup = server.Group("product-inventory-transactions", middleware.Authorize, middleware.AdminAuthorization)
	adminAuthGroup.GET("/product/:id", handler.GetInventoryTransactionsByProduct)
	adminAuthGroup.GET("/:id", handler.GetProductInventoryTransaction)
	adminAuthGroup.POST("/create", handler.CreateProductInventoryTransaction)
	adminAuthGroup.PUT("/update", handler.UpdateProductInventoryTransaction)
}
