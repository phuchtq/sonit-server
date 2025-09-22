package handler

import (
	action_type "sonit_server/constant/action_type"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	business_logic "sonit_server/usecase/business_logic"
	"sonit_server/utils"

	"github.com/gin-gonic/gin"
)

// ActivateCollection godoc
// @Summary Activate a collection
// @Description Activate a collection by ID
// @Tags collections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Collection ID"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /collections/activate/{id} [patch]
func ActivateCollection(ctx *gin.Context) {
	service, err := business_logic.GenerateCollectionService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:  service.ActivateCollection(ctx.Param("id"), ctx),
		Context: ctx,
	})
}

// CreateCollection godoc
// @Summary Create a new collection
// @Description Create a new collection with the provided details
// @Tags collections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param collection body request.CreateCollectionRequest true "Collection creation request"
// @Success 201 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /collections/create [post]
func CreateCollection(ctx *gin.Context) {
	var request request.CreateCollectionRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil)
		return
	}

	service, err := business_logic.GenerateCollectionService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:  service.CreateCollection(request, ctx),
		Context: ctx,
	})
}

// GetAllCollections godoc
// @Summary Get all collections
// @Description Retrieve all collections from the system
// @Tags collections
// @Accept json
// @Produce json
// @Success 200 {object} []entity.Collection
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /collections [get]
func GetAllCollections(ctx *gin.Context) {
	service, err := business_logic.GenerateCollectionService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	res, err := service.GetAllCollections(ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		PostType: action_type.NON_POST,
		Context:  ctx,
	})
}

// GetCollectionById godoc
// @Summary Get collection by ID
// @Description Retrieve a specific collection by its ID
// @Tags collections
// @Accept json
// @Produce json
// @Param id path string true "Collection ID"
// @Success 200 {object} entity.Collection
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /collections/{id} [get]
func GetCollectionById(ctx *gin.Context) {
	service, err := business_logic.GenerateCollectionService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	res, err := service.GetCollectionById(ctx.Param("id"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		PostType: action_type.NON_POST,
		Context:  ctx,
	})
}

// GetCollectionsByName godoc
// @Summary Get collections by name
// @Description Retrieve collections that match the specified name
// @Tags collections
// @Accept json
// @Produce json
// @Param name path string true "Collection name"
// @Success 200 {object} []entity.Collection
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /collections/name/{name} [get]
func GetCollectionsByName(ctx *gin.Context) {
	service, err := business_logic.GenerateCollectionService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	res, err := service.GetCollectionsByName(ctx.Param("name"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		PostType: action_type.NON_POST,
		Context:  ctx,
	})
}

// GetCollectionsByStatus godoc
// @Summary Get collections by status
// @Description Retrieve collections filtered by their status (active/inactive)
// @Tags collections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status body bool true "Collection status (true for active, false for inactive)"
// @Success 200 {object} []entity.Collection
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /collections/status [get]
func GetCollectionsByStatus(ctx *gin.Context) {
	var status bool
	if ctx.ShouldBindJSON(&status) != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil)
		return
	}

	service, err := business_logic.GenerateCollectionService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	res, err := service.GetCollectionsByStatus(&status, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		PostType: action_type.NON_POST,
		Context:  ctx,
	})
}

// RemoveCollection godoc
// @Summary Remove a collection
// @Description Remove/delete a collection by ID
// @Tags collections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Collection ID"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /collections/remove/{id} [delete]
func RemoveCollection(ctx *gin.Context) {
	service, err := business_logic.GenerateCollectionService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:  service.RemoveCollection(ctx.Param("id"), ctx),
		Context: ctx,
	})
}

// UpdateCollection godoc
// @Summary Update a collection
// @Description Update an existing collection with new details
// @Tags collections
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param collection body request.UpdateCollectionRequest true "Collection update request"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /collections/update [put]
func UpdateCollection(ctx *gin.Context) {
	var request request.UpdateCollectionRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil)
		return
	}

	service, err := business_logic.GenerateCollectionService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:  service.UpdateCollection(request, ctx),
		Context: ctx,
	})
}
