package dataaccess

import (
	"context"
	"sonit_server/model/dto/request"
	"sonit_server/model/entity"
)

type IShippingRepo interface {
	GetShippings(req request.GetShippingsRequest, ctx context.Context) (*[]entity.Shipping, int, error)
	GetShipping(id string, ctx context.Context) (*entity.Shipping, int, error)
	CreateShipping(ship entity.Shipping, ctx context.Context) error
	UpdateShipping(ship entity.Shipping, ctx context.Context) error
	RemoveShipping(id string, ctx context.Context) error
}
