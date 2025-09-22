package handler

import (
	action_type "sonit_server/constant/action_type"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	business_logic "sonit_server/usecase/business_logic"
	"sonit_server/utils"

	"github.com/gin-gonic/gin"
)

// ActivateCategory godoc
// @Summary Activate a category
// @Description Activate a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /categories/{id}/activate [patch]
func ActivateCategory(ctx *gin.Context) {
	service, err := business_logic.GenerateCategoryService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:  service.ActivateCategory(ctx.Param("id"), ctx),
		Context: ctx,
	})
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category with the provided details
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category body request.CreateCategoryRequest true "Category creation request"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /categories [post]
func CreateCategory(ctx *gin.Context) {
	var request request.CreateCategoryRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil)
		return
	}

	service, err := business_logic.GenerateCategoryService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:  service.CreateCategory(request, ctx),
		Context: ctx,
	})
}

// GetAllCategories godoc
// @Summary Get all categories
// @Description Retrieve all categories from the system
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []entity.Category
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /categories [get]
func GetAllCategories(ctx *gin.Context) {
	service, err := business_logic.GenerateCategoryService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	res, err := service.GetAllCategories(ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		PostType: action_type.NON_POST,
		Context:  ctx,
	})
}

// GetCategoryById godoc
// @Summary Get category by ID
// @Description Retrieve a specific category by its ID
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} []entity.Category
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /categories/{id} [get]
func GetCategoryById(ctx *gin.Context) {
	service, err := business_logic.GenerateCategoryService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	res, err := service.GetCategoryById(ctx.Param("id"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		PostType: action_type.NON_POST,
		Context:  ctx,
	})
}

// GetCategoriesByName godoc
// @Summary Get categories by name
// @Description Retrieve categories that match the specified name
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param name path string true "Category name"
// @Success 200 {object} []entity.Category
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /categories/name/{name} [get]
func GetCategoriesByName(ctx *gin.Context) {
	service, err := business_logic.GenerateCategoryService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	res, err := service.GetCategoriesByName(ctx.Param("name"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		PostType: action_type.NON_POST,
		Context:  ctx,
	})
}

// GetCategoriesByStatus godoc
// @Summary Get categories by status
// @Description Retrieve categories filtered by their status (active/inactive)
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status body bool true "Category status (true for active, false for inactive)"
// @Success 200 {object} []entity.Category
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /categories/status [get]
func GetCategoriesByStatus(ctx *gin.Context) {
	var status bool
	if ctx.ShouldBindJSON(&status) != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil)
		return
	}

	service, err := business_logic.GenerateCategoryService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	res, err := service.GetCategoriesByStatus(&status, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		PostType: action_type.NON_POST,
		Context:  ctx,
	})
}

// RemoveCategory godoc
// @Summary Remove a category
// @Description Remove/delete a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /categories/remove/{id} [delete]
func RemoveCategory(ctx *gin.Context) {
	service, err := business_logic.GenerateCategoryService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:  service.RemoveCategory(ctx.Param("id"), ctx),
		Context: ctx,
	})
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category with new details
// @Tags categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category body request.UpdateCategoryRequest true "Category update request"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /categories/update [put]
func UpdateCategory(ctx *gin.Context) {
	var request request.UpdateCategoryRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil)
		return
	}

	service, err := business_logic.GenerateCategoryService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:  service.UpdateCategory(request, ctx),
		Context: ctx,
	})
}
