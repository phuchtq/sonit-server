package main

import "sonit_server/cmd"

// @title Sonit Server API
// @version 1.0
// @description API for Sonit Server
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Execute server application
	cmd.Execute()
}
