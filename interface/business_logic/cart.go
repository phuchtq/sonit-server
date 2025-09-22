package businesslogic

import (
	"context"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
)

type ICartService interface {
	ViewCartDetail(id string, pageNumber int, ctx context.Context) (response.PaginationDataResponse, error)
	AddItemToCart(req request.AddItemToCartRequest, ctx context.Context) error
	RemoveItem(req request.RemoveItemFromCartRequest, ctx context.Context) error
	EditItemInCart(req request.EditItemInCartRequest, ctx context.Context) error
}
