package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	api_route "sonit_server/api_route"
	"sonit_server/constant/env"
	payment_env "sonit_server/constant/env/payment"
	"sonit_server/constant/noti"
	payment_method "sonit_server/constant/payment_method"
	"sonit_server/docs"
	_ "sonit_server/docs"

	"github.com/gin-gonic/gin"
	"github.com/payOSHQ/payos-lib-golang"
	swagger_files "github.com/swaggo/files"
	gin_swagger "github.com/swaggo/gin-swagger"
)

func setupApiRoutes(server *gin.Engine, port string) {
	// Role API endpoints
	api_route.InitializeRoleHandlerRoute(server, port)

	// User API endpoints
	api_route.InitializeUserHandlerRoute(server, port)

	// Collection API endpoints
	api_route.InitializeCollectionHandlerRoute(server, port)

	// Category API endpoints
	api_route.InitializeCategoryHandlerRoute(server, port)

	// Voucher API endpoints
	api_route.InitializeVoucherHandlerRoute(server, port)

	// Product API endpoints
	api_route.InitializeProductHandlerRoute(server, port)

	// Product Inventory Transaction API endpoints
	api_route.InitializeProductInventoryTransactionHandlerRoute(server, port)

	// Cart API endpoints
	api_route.InitializeCartHandlerRoute(server, port)

	// Payment API endpoints
	api_route.InitializePaymentHandlerRoute(server, port)

	// Order API endpoints
	api_route.InitializeOrderHandlerRoute(server, port)

	server.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusMovedPermanently, "/swagger/index.html#")
	})
}

func setupSwagger(server *gin.Engine, port string) {
	// Configure swagger info
	docs.SwaggerInfo.Title = "Sonit Server API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http"}
	docs.SwaggerInfo.Host = "localhost:" + port
	docs.SwaggerInfo.Host = os.Getenv(env.SWAGGER_HOST)

	// Add swagger route
	server.GET("swagger/*any", gin_swagger.WrapHandler(swagger_files.Handler))
}

func setupPayments(logger *log.Logger) {
	// Payos
	if err := payos.Key(os.Getenv(payment_env.PAYOS_CLIENT_ID), os.Getenv(payment_env.PAYOS_API_KEY), os.Getenv(payment_env.PAYOS_CHECKSUM_KEY)); err != nil {
		logger.Println(fmt.Sprintf(noti.PAYMENT_INIT_ENV_ERR_MSG, payment_method.PAYOS) + err.Error())
	}
}
