package dataaccess

import (
	"context"
	"sonit_server/model/dto/request"
	"sonit_server/model/entity"
)

type IProductInventoryTransactionRepo interface {
	GetAllProductInventoryTransactions(req request.GetProductInventoryTrasactionsRequest, ctx context.Context) (*[]entity.ProductInventoryTransaction, int, error)
	GetInventoryTransactionsByProduct(req request.GetProductInventoryTrasactionsByProductRequest, ctx context.Context) (*[]entity.ProductInventoryTransaction, int, error)
	GetProductInventoryTransaction(id string, ctx context.Context) (*entity.ProductInventoryTransaction, error)
	CreateProductInventoryTransaction(tx entity.ProductInventoryTransaction, ctx context.Context) error
	UpdateProductInventoryTransaction(tx entity.ProductInventoryTransaction, ctx context.Context) error
}
