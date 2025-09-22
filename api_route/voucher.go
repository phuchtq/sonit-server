package apiroute

import (
	"sonit_server/handler"
	"sonit_server/utils/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeVoucherHandlerRoute(server *gin.Engine, port string) {
	// Context path
	var contextPath string = "vouchers"

	// Define Voucher endpoints with admin required
	var adminAuthGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorization)
	adminAuthGroup.POST("", handler.CreateVoucher)
	adminAuthGroup.GET("", handler.GetAllVouchers)
	adminAuthGroup.GET("/valid", handler.GetAllValidVouchers)
	adminAuthGroup.PUT("/update", handler.UpdateVoucher)
	adminAuthGroup.DELETE("/remove/:id", handler.RemoveVoucher)

	// Define Voucher endpoints with basic required
	var norGroup = server.Group(contextPath)
	norGroup.GET("/:id", handler.GetVoucherById)
}
