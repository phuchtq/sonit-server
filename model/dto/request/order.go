package request

import "sonit_server/model/dto/response"

type GetOrdersRequest struct {
	Request SearchPaginatioRequest `json:"request"`
	UserId  string                 `json:"user_id" form:"user_id"`
	Status  string                 `json:"status" form:"status"`
}

type CreateOrderRequest struct {
	UserId       string              `json:"user_id" validate:"required"`
	ShippingUnit string              `json:"shipping_unit"`
	Items        []response.CartItem `json:"items" validate:"required, min=1"`
	Status       string              `json:"status"`
	Note         string              `json:"note"`
}

type UpdateOrderRequest struct {
	OrderId  string `json:"order_id" validate:"required"`
	Currency string `json:"currency"`
	Status   string `json:"status"`
	Note     string `json:"note"`
}
