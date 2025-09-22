package cmd

import (
	"os"
	"sonit_server/constant/env"
	"sonit_server/utils"

	"github.com/gin-gonic/gin"
)

// Execute server application
func Execute() {
	// Initialize logger config
	var logger = utils.GetLogConfig()

	// Load env
	loadEnv(logger)

	// Initialize gin server for API
	var server = gin.Default()

	// Config CORS for requests
	corsConfig(server)

	// Get API port
	var apiPort = os.Getenv(env.API_PORT)

	// Set up API routes
	setupApiRoutes(server, apiPort)

	// Set up swagger
	setupSwagger(server, apiPort)

	// Setup payments
	setupPayments(logger)

	// Run server
	if err := server.Run(":" + apiPort); err != nil {
		logger.Println("Error run server - " + err.Error())
	}
}
