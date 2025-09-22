package handler

import (
	action_type "sonit_server/constant/action_type"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	business_logic "sonit_server/usecase/business_logic"
	"sonit_server/utils"

	"github.com/gin-gonic/gin"
)

// GetAllPayments godoc
// @Summary Get all payments
// @Description Retrieve a paginated list of payments
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param        page query int false "Page number"
// @Param        status query string false "Payment status"
// @Param        method query string false "Payment method"
// @Param        user_id query string false "User ID"
// @Success 200 {object} response.PaginationDataResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /payments [get]
func GetAllPayments(ctx *gin.Context) {
	var request request.GetPaymentsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GeneratePaymentService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetPayments(request, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetPaymentById godoc
// @Summary Get payment by ID
// @Description Retrieve a single payment record by its ID
// @Tags payments
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payment ID"
// @Success 200 {object} entity.Payment
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /payments/{id} [get]
func GetPaymentById(ctx *gin.Context) {
	service, err := business_logic.GeneratePaymentService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.GetPaymentById(ctx.Param("id"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// GetPaymentsByUser godoc
// @Summary Get payments by user ID
// @Description Retrieve a list of payments made by a specific user
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param        id path string true "User ID"
// @Param        page_number query int false "Page number"
// @Param        keyword     query string false "Search keyword"
// @Param        filter_prop query string false "Filter property (e.g. date, price)"
// @Param        order       query string false "Sort order (ASC or DESC)"
// @Param        status       query string false "Payment status (e.g. pending, paid)"
// @Param        method       query string false "Payment Method (e.g. VTP, GHTK)"
// @Success 200 {object} response.PaginationDataResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /payments/user/{id} [get]
func GetPaymentsByUser(ctx *gin.Context) {
	var request request.GetPaymentsRequest
	if ctx.ShouldBindQuery(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GeneratePaymentService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	request.UserId = ctx.Param("id")

	res, err := service.GetPayments(request, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// UpdatePayment godoc
// @Summary Update a payment record
// @Description Update payment details
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.UpdatePaymentRequest true "UpdatePaymentRequest"
// @Success 200 {object} response.MessageAPIResponse "Success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /payments/update [put]
func UpdatePayment(ctx *gin.Context) {
	var request request.UpdatePaymentRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GeneratePaymentService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:   service.UpdatePayment(request, ctx),
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// CreatePayment godoc
// @Summary      Create a payment via cart
// @Description  Creates a new payment through cart and returns redirect info or confirmation
// @Tags         payments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body request.CreatePaymentThroughCartRequest true "Create Payment Through Cart Request"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /payments/create/cart [post]
func CreatePaymentThroughCart(ctx *gin.Context) {
	var request request.CreatePaymentThroughCartRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GeneratePaymentService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreatePaymentThroughCart(request, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}

// CreatePaymentDirect godoc
// @Summary      Create direct payment
// @Description  Creates a payment and processes it directly without redirect
// @Tags         payments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body request.CreatePaymentDirectRequest true "Direct Payment Request"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /payments/create/direct [post]
func CreatePaymentDirect(ctx *gin.Context) {
	var request request.CreatePaymentDirectRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GeneratePaymentService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CreatePaymentDirect(request, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}

// CallbackPaymentSuccess godoc
// @Summary      Callback after successful payment
// @Description  Handles redirect or callback after successful payment
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        id path string true "Payment ID"
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Router       /payments/callback-success/{id} [get]
func CallbackPaymentSuccess(ctx *gin.Context) {
	service, err := business_logic.GeneratePaymentService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CallbackPaymentSuccess(ctx.Param("id"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.REDIRECT,
	})
}

// CallbackPaymentCancel godoc
// @Summary      Callback after canceled payment
// @Description  Handles redirect or callback after canceled payment
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        id path string true "Payment ID"
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Router       /payments/callback-cancel/{id} [get]
func CallbackPaymentCancel(ctx *gin.Context) {
	service, err := business_logic.GeneratePaymentService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	res, err := service.CallbackPaymentCancel(ctx.Param("id"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		Data2:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.REDIRECT,
	})
}
