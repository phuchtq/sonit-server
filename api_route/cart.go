package apiroute

import (
	"sonit_server/handler"
	"sonit_server/utils/middleware"

	"github.com/gin-gonic/gin"
)

func InitializeCartHandlerRoute(server *gin.Engine, port string) {
	var authGroup = server.Group("carts", middleware.Authorize)
	authGroup.GET("/:id", handler.ViewCartDetail)
	authGroup.POST("/item/add", handler.AddItemToCart)
	authGroup.PUT("/item/edit", handler.EditItemInCart)
	authGroup.DELETE("/item/remove", handler.RemoveItemInCart)
}
