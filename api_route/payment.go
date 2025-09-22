package apiroute

import (
	"sonit_server/handler"
	"sonit_server/utils/middleware"

	"github.com/gin-gonic/gin"
)

func InitializePaymentHandlerRoute(server *gin.Engine, port string) {
	// Context path
	var contextPath string = "payments"

	// Define Payment endpoints with admin required
	var adminAuthGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorization)
	adminAuthGroup.GET("", handler.GetAllPayments)
	adminAuthGroup.PUT("/update", handler.UpdatePayment)

	// Define Payment endpoints with basic required
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.GET("/user/:id", handler.GetPaymentsByUser)
	authGroup.GET("/:id", handler.GetPaymentById)
	authGroup.POST("/create/cart", handler.CreatePaymentThroughCart)
	authGroup.POST("/create/direct", handler.CreatePaymentDirect)

	var norGroup = server.Group(contextPath)
	norGroup.GET("/callback-success/:id", handler.CallbackPaymentSuccess)
	norGroup.GET("/callback-cancel/:id", handler.CallbackPaymentCancel)
}
