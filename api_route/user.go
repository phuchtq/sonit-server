package apiroute

import (
	"sonit_server/handler"
	"sonit_server/utils/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeUserHandlerRoute(server *gin.Engine, port string) {
	// Context path
	var contextPath string = "users"

	var adminAuthGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorization)
	adminAuthGroup.GET("", handler.GetAllUsers)
	adminAuthGroup.GET("/role/:role", handler.GetUsersByRole)
	adminAuthGroup.GET("/status/:status", handler.GetUsersByStatus)

	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.GET("/:id", handler.GetUser)
	authGroup.PUT("/update-info/:actorId", handler.UpdateUser)
	authGroup.POST("/vip/create", handler.CreateVipCode)
	//authGroup.PUT("/id/:id/status/:status", handler.ChangeUserStatus)
	//authGroup.PUT("/logout/:userId", handler.Logout)

	var norGroup = server.Group(contextPath)
	//norGroup.POST("/login", handler.Login)
	norGroup.POST("/create", handler.CreateUser)
	//norGroup.PUT("/:email", handler.Re)
	norGroup.PUT("/reset-password/:token/:password/:confirmPassword")
	norGroup.GET("/verify-action", handler.VerifyAction)

	// Auth group with login, logout
	var norCredentialGroup = server.Group("auth")
	norCredentialGroup.POST("/login", handler.Login)
	norCredentialGroup.POST("/refresh-token", handler.RefreshToken)

	var authCredentialGroup = server.Group("auth", middleware.Authorize)
	authCredentialGroup.POST("/logout/:id", handler.Logout)
}
