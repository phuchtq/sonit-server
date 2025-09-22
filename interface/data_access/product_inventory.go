package dataaccess

import (
	"context"
	"sonit_server/model/entity"
)

type IProductInventoryRepo interface {
	GetProductInventory(id string, ctx context.Context) (*entity.ProductInventory, error)
	UpdateProductInventory(inventory entity.ProductInventory, ctx context.Context) error
	CreateProductInventory(inventory entity.ProductInventory, ctx context.Context) error
}
