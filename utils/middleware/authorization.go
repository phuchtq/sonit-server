package middleware

import (
	"log"
	"os"
	"sonit_server/constant/env"
	"sonit_server/utils"

	"github.com/gin-gonic/gin"
)

func Authorize(ctx *gin.Context) {
	// Get token from the header
	var token string = ctx.Request.Header.Get("Authorization")
	log.Println("Token: ", token)

	var unAuthBodyResponse = utils.GetUnAuthBodyResponse(ctx)

	if token == "" {
		utils.ProcessResponse(unAuthBodyResponse)
		ctx.Abort()
		return
	}

	userId, role, _, err := utils.ExtractDataFromToken(token, utils.GetLogConfig())
	if err != nil {
		utils.ProcessResponse(unAuthBodyResponse)
		ctx.Abort()
		return
	}

	log.Println("Role: ", role)

	ctx.Set("userId", userId)
	ctx.Set("role", role)
	ctx.Next()
}

func AdminAuthorization(ctx *gin.Context) {
	log.Println("Role access: ", ctx.GetString("role"))
	if ctx.GetString("role") != os.Getenv(env.ADMIN_ROLE) {
		utils.ProcessResponse(utils.GetUnAuthBodyResponse(ctx))
		ctx.Abort()
		return
	}

	ctx.Next()
}
