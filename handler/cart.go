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

// ViewCartDetail godoc
// @Summary      View cart detail
// @Description  Retrieves the details of a user's cart
// @Tags         carts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id         path      string  true  "User ID"
// @Param        pageNumber query     int     false "Page number for pagination"
// @Success      200 {object} response.PaginationDataResponse
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /carts/{id} [get]
func ViewCartDetail(ctx *gin.Context) {
	service, err := business_logic.GenerateCartService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	pageNumber, _ := strconv.Atoi(ctx.Query("pageNumber"))

	res, err := service.ViewCartDetail(ctx.Param("id"), pageNumber, ctx)

	utils.ProcessResponse(response.APIResponse{
		Data1:    res,
		ErrMsg:   err,
		Context:  ctx,
		PostType: action_type.NON_POST,
	})
}

// AddItemToCart godoc
// @Summary      Add item to cart
// @Description  Adds a new item to the user's cart
// @Tags         carts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body     request.AddItemToCartRequest true "Item to add"
// @Success 200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /carts/item/add [post]
func AddItemToCart(ctx *gin.Context) {
	var request request.AddItemToCartRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GenerateCartService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:   service.AddItemToCart(request, ctx),
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}

// RemoveItemInCart godoc
// @Summary      Remove item from cart
// @Description  Removes an item from the user's cart
// @Tags         carts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body     request.RemoveItemFromCartRequest true "Item to remove"
// @Success      200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /carts/item/remove [delete]
func RemoveItemInCart(ctx *gin.Context) {
	var request request.RemoveItemFromCartRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GenerateCartService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:   service.RemoveItem(request, ctx),
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}

// EditItemInCart godoc
// @Summary      Edit item in cart
// @Description  Updates the quantity or details of an item in the cart
// @Tags         carts
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body     request.EditItemInCartRequest true "Item to edit"
// @Success      200 {object} response.MessageAPIResponse "success"
// @Failure 401 {object} response.MessageAPIResponse "You have no rights to access this action."
// @Failure 400 {object} response.MessageAPIResponse "Invalid data. Please try again."
// @Failure 500 {object} response.MessageAPIResponse "There is something wrong in the system during the process. Please try again later."
// @Router       /carts/item/edit [put]
func EditItemInCart(ctx *gin.Context) {
	var request request.EditItemInCartRequest
	if ctx.ShouldBindJSON(&request) != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, nil))
		return
	}

	service, err := business_logic.GenerateCartService()
	if err != nil {
		utils.ProcessResponse(utils.GenerateInvalidRequestAndSystemProblemModel(ctx, err))
		return
	}

	utils.ProcessResponse(response.APIResponse{
		ErrMsg:   service.EditItemInCart(request, ctx),
		Context:  ctx,
		PostType: action_type.INFORM,
	})
}
