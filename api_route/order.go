package apiroute

import (
	"sonit_server/handler"
	"sonit_server/utils/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeOrderHandlerRoute(server *gin.Engine, port string) {
	// Context path
	var contextPath string = "orders"

	// Define Order endpoints with admin required
	var adminAuthGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorization)
	adminAuthGroup.GET("", handler.GetOrders)

	// Define Order endpoints with basic required
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.GET("/user/:id", handler.GetOrdersByUser)
	authGroup.GET("/:id", handler.GetOrder)
	authGroup.POST("/create", handler.CreateOrder)
	authGroup.PUT("/update", handler.UpdateOrder)
}
