package dataaccess

import (
	"context"
	"sonit_server/model/dto/request"
	"sonit_server/model/entity"
)

type IOrderRepo interface {
	GetOrders(req request.GetOrdersRequest, ctx context.Context) (*[]entity.Order, int, error)
	GetOrder(id string, ctx context.Context) (*entity.Order, error)
	CreateOrder(order entity.Order, ctx context.Context) error
	UpdateOrder(order entity.Order, ctx context.Context) error
}
