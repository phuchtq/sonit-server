package businesslogic

import (
	"context"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
)

type IOrderService interface {
	GetOrders(req request.GetOrdersRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetOrder(id string, ctx context.Context) (*response.ViewOrderResponse, error)
	CreateOrder(req request.CreateOrderRequest, ctx context.Context) error
	UpdateOrder(req request.UpdateOrderRequest, ctx context.Context) error
}
