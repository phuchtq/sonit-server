package apiroute

import (
	"sonit_server/handler"
	"sonit_server/utils/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeCollectionHandlerRoute(server *gin.Engine, port string) {
	// Context path
	var contextPath string = "collections"

	// Define Collection endpoints with admin required
	var adminAuthGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorization)
	adminAuthGroup.GET("/status", handler.GetCollectionsByStatus)
	adminAuthGroup.POST("/create", handler.CreateCollection)
	adminAuthGroup.PUT("/update", handler.UpdateCollection)
	adminAuthGroup.PATCH("activate/:id", handler.ActivateCollection)
	adminAuthGroup.DELETE("remove/:id", handler.RemoveCollection)

	// Define Collection endpoints with basic required
	var norGroup = server.Group(contextPath)
	norGroup.GET("", handler.GetAllCollections)
	norGroup.GET("/:id", handler.GetCollectionById)
}
