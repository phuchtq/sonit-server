package apiroute

import (
	"sonit_server/handler"
	"sonit_server/utils/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeCategoryHandlerRoute(server *gin.Engine, port string) {
	// Context path
	var contextPath string = "categories"

	// Define Category endpoints with admin required
	var adminAuthGroup = server.Group(contextPath, middleware.Authorize, middleware.AdminAuthorization)
	adminAuthGroup.GET("", handler.GetAllCategories)
	adminAuthGroup.GET("/name/:name", handler.GetCategoriesByName)
	adminAuthGroup.GET("/status", handler.GetCategoriesByStatus)
	adminAuthGroup.POST("/create", handler.CreateCategory)
	adminAuthGroup.PUT("/update", handler.UpdateCategory)
	adminAuthGroup.PATCH("/activate/:id", handler.ActivateCategory)
	adminAuthGroup.DELETE("remove/:id", handler.RemoveCategory)

	// Define Category endpoints with basic required
	var norGroup = server.Group(contextPath)
	norGroup.GET("/:id", handler.GetCategoryById)
}
