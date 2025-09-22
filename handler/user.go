package handler

import (
	action_type "sonit_server/constant/action_type"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	business_logic "sonit_server/usecase/business_logic"
	"sonit_server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Summary      User login
// @Description  Authenticates user credentials and returns access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body request.LoginRequest true "Login credentials"
// @Success      200 {object} response.LoginSuccessResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /auth/login [post]
func Login(ctx *gin.Context) {
	var request request.LoginRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GenerateUserService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res1, res2, err := service.Login(request, ctx)

	utils.ProcessLoginResponse(response.APIResponse{
		Data1:   res1,
		Data2:   res2,
		ErrMsg:  err,
		Context: ctx,
	})
}

// Logout godoc
// @Summary      Logout user
// @Description  Logs the user out and invalidates the session
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /auth/logout/{id} [post]
func Logout(ctx *gin.Context) {
	service, err := business_logic.GenerateUserService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:   service.Logout(ctx.Param("id"), ctx),
		PostType: action_type.INFORM,
		Context:  ctx,
	})
}

// GetUsersByRole godoc
// @Summary      Get users by role
// @Description  Retrieves a list of users assigned to a specific role
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        role       path  string true  "Role name"
// @Param        pageNumber query int    false "Page number"
// @Success      200 {object} response.PaginationDataResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /users/role/{role} [get]
func GetUsersByRole(ctx *gin.Context) {
	service, err := business_logic.GenerateUserService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	pageNumber, _ := strconv.Atoi(ctx.Query("pageNumber"))

	res, err := service.GetUsersByRole(pageNumber, ctx.Param("role"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetUser godoc
// @Summary      Get user details
// @Description  Fetches a specific user's details
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} entity.User
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /users/{id} [get]
func GetUser(ctx *gin.Context) {
	service, err := business_logic.GenerateUserService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetUser(ctx.Param("id"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Retrieves all users with pagination
// @Tags         users
// @Security     BearerAuth
// @Produce      json
// @Param        pageNumber query int false "Page number"
// @Success      200 {object} response.PaginationDataResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /users [get]
func GetAllUsers(ctx *gin.Context) {
	service, err := business_logic.GenerateUserService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	pageNumber, _ := strconv.Atoi(ctx.Query("pageNumber"))

	res, err := service.GetAllUsers(pageNumber, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetUsersByStatus godoc
// @Summary      Get users by status
// @Description  Returns users filtered by active/inactive status
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        pageNumber query int false "Page number"
// @Param        status path string false "Status filter (true for active, false for inactive)"
// @Success      200 {object} response.PaginationDataResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /users/status/ [get]
func GetUsersByStatus(ctx *gin.Context) {
	service, err := business_logic.GenerateUserService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	pageNumber, _ := strconv.Atoi(ctx.Query("pageNumber"))

	var pStatus *bool
	status, err := strconv.ParseBool(ctx.Param("status"))
	if err == nil {
		pStatus = &status
	}

	res, err := service.GetUsersByStatus(pageNumber, pStatus, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// UpdateUser godoc
// @Summary      Update user
// @Description  Updates an existing user's information
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body request.UpdateUserRequest true "User update details"
// @Success      200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /users/update-info [put]
func UpdateUser(ctx *gin.Context) {
	var request request.UpdateUserRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil)
		return
	}

	service, err := business_logic.GenerateUserService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.UpdateUser(request, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}

// VerifyAction godoc
// @Summary      Verify action
// @Description  Verifies user action with a raw token (email, password reset, etc.)
// @Tags         users
// @Produce      json
// @Param        rawToken query string true "Raw verification token"
// @Success      200 {object} interface{}
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /users/verify-action [get]
func VerifyAction(ctx *gin.Context) {
	service, err := business_logic.GenerateUserService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	var rawToken = ctx.Query("rawToken")
	res, err := service.VerifyAction(rawToken, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.REDIRECT,
	})
}

// ResetPassword godoc
// @Summary      Reset password
// @Description  Resets user password using a token
// @Tags         users
// @Produce      json
// @Param        password        path string true "New password"
// @Param        confirmPassword path string true "Confirm password"
// @Param        token         path string true "Reset token"
// @Success      200 {object} interface{}
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /users/reset-password/{token}/{password}/{confirmPassword} [post]
func ResetPassword(ctx *gin.Context) {
	service, err := business_logic.GenerateUserService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.ResetPassword(ctx.Param("password"), ctx.Param("confirmPassword"), ctx.Param("token"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.REDIRECT,
	})
}

// CreateUser godoc
// @Summary      Create new user
// @Description  Registers a new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body request.CreateUserRequest true "New user details"
// @Success      200 {object} response.MessageAPIResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /users/create [post]
func CreateUser(ctx *gin.Context) {
	var request request.CreateUserRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil)
		return
	}

	service, err := business_logic.GenerateUserService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreateUser(request, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}

// CreateVipCode godoc
// @Summary      Create vip code
// @Description  Registers vip code for vip user
// @Tags         users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body request.CreateVipCodeRequest true "Vip code-create details"
// @Success      200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /users/vip/create [post]
func CreateVipCode(ctx *gin.Context) {
	var request request.CreateVipCodeRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil)
		return
	}

	service, err := business_logic.GenerateUserService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:   service.CreateVipCode(request, ctx),
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Refreshes the access token using a valid refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body request.RefreshTokenRequest true "Refresh Token Request"
// @Success      200 {object} response.RefreshTokenSuccessResponse
// @Failure      401 {object} response.MessageAPIResponse
// @Router       /auth/refresh-token [post]
func RefreshToken(ctx *gin.Context) {
	var request request.RefreshTokenRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil)
		return
	}

	service, err := business_logic.GenerateUserService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.RefreshToken(request, ctx)

	utils.ProcessRefreshTokenResponse(response.APIResponse{
		Data1:   res,
		ErrMsg:  err,
		Context: ctx,
	})
}
