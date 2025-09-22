package handler

import (
	action_type "sonit_server/constant/action_type"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	business_logic "sonit_server/usecase/business_logic"
	"sonit_server/utils"

	"github.com/gin-gonic/gin"
)

// GetOrders godoc
// @Summary      Get all orders
// @Description  Retrieve a paginated list of orders
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page_number query int false "Page number"
// @Param        keyword     query string false "Search keyword"
// @Param        filter_prop query string false "Filter property (e.g. date, price)"
// @Param        order       query string false "Sort order (ASC or DESC)"
// @Param        user_id     query string false "User ID"
// @Success      200 {object} response.PaginationDataResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /orders [get]
func GetOrders(ctx *gin.Context) {
	var request request.GetOrdersRequest
	if ctx.ShouldBindQuery(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GenerateOrderService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetOrders(request, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetOrdersByUser godoc
// @Summary      Get orders by user ID
// @Description  Retrieve a list of orders placed by a specific user
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Param        page_number query int false "Page number"
// @Param        keyword     query string false "Search keyword"
// @Param        filter_prop query string false "Filter property (e.g. date, price)"
// @Param        order       query string false "Sort order (ASC or DESC)"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /orders/user/{id} [get]
func GetOrdersByUser(ctx *gin.Context) {
	var req request.SearchPaginatioRequest
	if ctx.ShouldBindQuery(&req) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GenerateOrderService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetOrders(request.GetOrdersRequest{
		Request: req,
		UserId:  ctx.Param("id"),
	}, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetOrder godoc
// @Summary      Get order by ID
// @Description  Retrieve a specific order by its ID
// @Tags         orders
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Order ID"
// @Success      200 {object} response.ViewOrderResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /orders/{id} [get]
func GetOrder(ctx *gin.Context) {
	service, err := business_logic.GenerateOrderService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetOrder(ctx.Param("id"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreateOrder godoc
// @Summary      Create a new order
// @Description  Submit a new order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body request.CreateOrderRequest true "CreateOrderRequest"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /orders/create [put]
func CreateOrder(ctx *gin.Context) {
	var request request.CreateOrderRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GenerateOrderService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:   service.CreateOrder(request, ctx),
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}

// UpdateOrder godoc
// @Summary      Update an existing order
// @Description  Modify an order by providing updated information
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body request.UpdateOrderRequest true "UpdateOrderRequest"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /orders/update [put]
func UpdateOrder(ctx *gin.Context) {
	var request request.UpdateOrderRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GenerateOrderService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:   service.UpdateOrder(request, ctx),
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}
