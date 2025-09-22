package handler

import (
	action_type "sonit_server/constant/action_type"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	business_logic "sonit_server/usecase/business_logic"
	"sonit_server/utils"

	"github.com/gin-gonic/gin"
)

// @Summary Activate a role
// @Description Activates a role by ID
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /roles/activate/{id} [patch]
func ActivateRole(ctx *gin.Context) {
	service, err := business_logic.GenerateRoleService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	var id = ctx.Param("id")

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:  service.ActivateRole(id, ctx),
		Context: ctx,
	})
}

// CreateRole creates a new role
// @Summary Create a new role
// @Description Create a new role with the specified name
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name path string true "Role Name"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /roles/create/{name} [post]
func CreateRole(ctx *gin.Context) {
	service, err := business_logic.GenerateRoleService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	var name = ctx.Param("name")

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:  service.CreateRole(name, ctx),
		Context: ctx,
	})
}

// GetAllRoles retrieves all roles
// @Summary Get all roles
// @Description Retrieve a list of all roles in the system
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []entity.Role
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /roles [get]
func GetAllRoles(ctx *gin.Context) {
	service, err := business_logic.GenerateRoleService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	res, err := service.GetAllRoles(ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		PostType: action_type.NON_POST,
		Context:  ctx,
	})
}

// GetRoleById retrieves a role by ID
// @Summary Get role by ID
// @Description Retrieve a specific role by its unique identifier
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 200 {object} entity.Role
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /roles/{id} [get]
func GetRoleById(ctx *gin.Context) {
	service, err := business_logic.GenerateRoleService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	var id = ctx.Param("id")
	res, err := service.GetRoleById(id, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		PostType: action_type.NON_POST,
		Context:  ctx,
	})
}

// GetRolesByName retrieves roles by name
// @Summary Get roles by name
// @Description Retrieve roles that match the specified name
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name path string false "Role name"
// @Success 200 {object} []entity.Role
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /roles/name/{name} [get]
func GetRolesByName(ctx *gin.Context) {
	service, err := business_logic.GenerateRoleService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	var name = ctx.Param("name")
	res, err := service.GetRolesByName(name, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		PostType: action_type.NON_POST,
		Context:  ctx,
	})
}

// GetRolesByStatus retrieves roles by status
// @Summary Get roles by status
// @Description Retrieve roles that match the specified status (active/inactive)
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status body bool false "Role status (true for active, false for inactive)"
// @Success 200 {object} []entity.Role
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /roles/status [get]
func GetRolesByStatus(ctx *gin.Context) {
	var status bool
	if ctx.ShouldBindJSON(&status) != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil)
		return
	}

	service, err := business_logic.GenerateRoleService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	res, err := service.GetRolesByStatus(&status, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		PostType: action_type.NON_POST,
		Context:  ctx,
	})
}

// RemoveRole removes a role
// @Summary Remove a role
// @Description Remove a role with the ID from the system
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name path string true "Role ID"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /roles/{name} [delete]
func RemoveRole(ctx *gin.Context) {
	service, err := business_logic.GenerateRoleService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:  service.RemoveRole(ctx.Param("id"), ctx),
		Context: ctx,
	})
}

// UpdateRole updates an existing role
// @Summary Update a role
// @Description Update an existing role with new information
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateRoleRequest true "Role update request"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /roles/update [put]
func UpdateRole(ctx *gin.Context) {
	var request request.UpdateRoleRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil)
		return
	}

	service, err := business_logic.GenerateRoleService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:  service.UpdateRole(request, ctx),
		Context: ctx,
	})
}
