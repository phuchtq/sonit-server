package handler

import (
	action_type "sonit_server/constant/action_type"
	request "sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	business_logic "sonit_server/usecase/business_logic"
	"sonit_server/utils"

	"github.com/gin-gonic/gin"
)

// GetAllVouchers godoc
// @Summary Get all vouchers
// @Description Retrieve all vouchers from the system
// @Tags vouchers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []entity.Voucher
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /vouchers [get]
func GetAllVouchers(ctx *gin.Context) {
	service, err := business_logic.GenerateVoucherService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	res, err := service.GetAllVouchers(ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		PostType: action_type.NON_POST,
		Context:  ctx,
	})
}

// GetAllValidVouchers godoc
// @Summary Get all valid vouchers
// @Description Retrieve all currently valid/active vouchers from the system
// @Tags vouchers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} []entity.Voucher
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /vouchers/valid [get]
func GetAllValidVouchers(ctx *gin.Context) {
	service, err := business_logic.GenerateVoucherService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	res, err := service.GetAllValidVouchers(ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		PostType: action_type.NON_POST,
		Context:  ctx,
	})
}

// GetVoucherById godoc
// @Summary Get voucher by ID
// @Description Retrieve a specific voucher by its ID
// @Tags vouchers
// @Accept json
// @Produce json
// @Param id path string true "Voucher ID"
// @Success 200 {object} []entity.Voucher
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /vouchers/{id} [get]
func GetVoucherById(ctx *gin.Context) {
	service, err := business_logic.GenerateVoucherService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	res, err := service.GetVoucherByID(ctx.Param("id"), ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		PostType: action_type.NON_POST,
		Context:  ctx,
	})
}

// CreateVoucher godoc
// @Summary Create a new voucher
// @Description Create a new voucher with the provided details
// @Tags vouchers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param voucher body request.CreateVoucherRequest true "Voucher creation request"
// @Success 201 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /vouchers/create [post]
func CreateVoucher(ctx *gin.Context) {
	var request request.CreateVoucherRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil)
		return
	}

	service, err := business_logic.GenerateVoucherService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:  service.CreateVoucher(request, ctx),
		Context: ctx,
	})
}

// UpdateVoucher godoc
// @Summary Update a voucher
// @Description Update an existing voucher with new details
// @Tags vouchers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param voucher body request.UpdateVoucherRequest true "Voucher update request"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /vouchers/update [put]
func UpdateVoucher(ctx *gin.Context) {
	var request request.UpdateVoucherRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil)
		return
	}

	service, err := business_logic.GenerateVoucherService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:  service.UpdateVoucher(request, ctx),
		Context: ctx,
	})
}

// RemoveVoucher godoc
// @Summary Remove a voucher
// @Description Remove/delete a voucher by ID
// @Tags vouchers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Voucher ID"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router /vouchers/remove/{id} [delete]
func RemoveVoucher(ctx *gin.Context) {
	service, err := business_logic.GenerateVoucherService()
	if err != nil {
		utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err)
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:  service.RemoveVoucher(ctx.Param("id"), ctx),
		Context: ctx,
	})
}
