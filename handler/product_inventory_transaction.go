package handler

import (
	action_type "sonit_server/constant/action_type"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	business_logic "sonit_server/usecase/business_logic"
	"sonit_server/utils"

	"github.com/gin-gonic/gin"
)

// @Summary Get inventory transactions by product
// @Description Retrieves inventory transactions filtered by product criteria
// @Tags product-inventory-transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param        id path string true "Product ID"
// @Param        page_number query int false "Page number"
// @Param        keyword     query string false "Search keyword"
// @Param        filter_prop query string false "Filter property (e.g. date, price)"
// @Param        order       query string false "Sort order (ASC or DESC)"
// @Param        action      query string false "Inventory action (e.g. import, export, sale, return)"
// @Success 200 {object} response.PaginationDataResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /product-inventory-transactions/product/{id} [get]
func GetInventoryTransactionsByProduct(ctx *gin.Context) {
	var req request.GetProductInventoryTrasactionsRequest
	if ctx.ShouldBindQuery(&req) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GenerateProductInventoryTransactionService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetInventoryTransactionsByProduct(request.GetProductInventoryTrasactionsByProductRequest{
		InventoryTransaction: req,
		ProductId:            ctx.Param("id"),
	}, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// @Summary Get product inventory transaction by ID
// @Description Retrieves a specific product inventory transaction by its ID
// @Tags product-inventory-transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product inventory transaction ID"
// @Success 200 {object} entity.ProductInventoryTransaction
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /product-inventory-transactions/{id} [get]
func GetProductInventoryTransaction(ctx *gin.Context) {
	service, err := business_logic.GenerateProductInventoryTransactionService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetProductInventoryTransaction(ctx.Param("id"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// @Summary Create product inventory transaction
// @Description Creates a new product inventory transaction record
// @Tags product-inventory-transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.CreateProductInventoryTransactionRequest true "Create product inventory transaction request"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /product-inventory-transactions/create [post]
func CreateProductInventoryTransaction(ctx *gin.Context) {
	var request request.CreateProductInventoryTransactionRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GenerateProductInventoryTransactionService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:   service.CreateProductInventoryTransaction(request, ctx),
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}

// @Summary Update product inventory transaction
// @Description Updates an existing product inventory transaction
// @Tags product-inventory-transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdateProductInventoryTransactionRequest true "Update product inventory transaction request"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Router /product-inventory-transactions/update [put]
func UpdateProductInventoryTransaction(ctx *gin.Context) {
	var request request.UpdateProductInventoryTransactionRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GenerateProductInventoryTransactionService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:   service.UpdateProductInventoryTransaction(request, ctx),
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}
