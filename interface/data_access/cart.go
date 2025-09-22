package dataaccess

import (
	"context"
	"sonit_server/model/entity"
)

type ICartRepo interface {
	GetCart(id string, ctx context.Context) (*entity.Cart, int, error)
	UpdateCart(cart entity.Cart, ctx context.Context) error
	CreateCart(cart entity.Cart, ctx context.Context) error
	RemoveCart(id string, ctx context.Context) error
}
