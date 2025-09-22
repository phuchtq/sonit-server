package businesslogic

import (
	"context"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	entity "sonit_server/model/entity"
)

type IProductInventoryTransactionService interface {
	GetAllProductInventoryTransactions(req request.GetProductInventoryTrasactionsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetInventoryTransactionsByProduct(req request.GetProductInventoryTrasactionsByProductRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetProductInventoryTransaction(id string, ctx context.Context) (*entity.ProductInventoryTransaction, error)
	CreateProductInventoryTransaction(req request.CreateProductInventoryTransactionRequest, ctx context.Context) error
	UpdateProductInventoryTransaction(req request.UpdateProductInventoryTransactionRequest, ctx context.Context) error
}
