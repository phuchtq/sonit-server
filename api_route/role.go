package apiroute

import (
	"sonit_server/handler"
	"sonit_server/utils/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeRoleHandlerRoute(server *gin.Engine, port string) {
	// Context path
	var contextPath string = "roles"

	// Define role endpoints with admin required
	var adminAuthGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorization)
	adminAuthGroup.GET("", handler.GetAllRoles)
	adminAuthGroup.GET("/name/:name", handler.GetRolesByName)
	adminAuthGroup.GET("/status", handler.GetRolesByStatus)
	adminAuthGroup.POST("create/:name", handler.CreateRole)
	adminAuthGroup.PUT("/update", handler.UpdateRole)
	adminAuthGroup.PATCH("activate/:id", handler.ActivateRole)
	adminAuthGroup.DELETE("remove/:id", handler.RemoveRole)

	// Define role endpoints with basic required
	var authGroup = server.Group(contextPath, middleware.Authorize)
	authGroup.GET("/:id", handler.GetRoleById)
}
