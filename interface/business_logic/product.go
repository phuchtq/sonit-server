package businesslogic

import (
	"context"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	entity "sonit_server/model/entity"
)

type IProductService interface {
	GetAllProducts(pageNumber int, ctx context.Context) (response.PaginationDataResponse, error)
	GetProductsCustomerUI(req request.GetProductsCustomerUI, ctx context.Context) (response.PaginationDataResponse, error)
	GetProductsByCategory(pageNumber int, id string, ctx context.Context) (response.PaginationDataResponse, error)
	GetProductsByCollection(pageNumber int, id string, ctx context.Context) (response.PaginationDataResponse, error)
	GetProductsByPriceInterval(pageNumber int, maxPrice, minPrice int64, ctx context.Context) (response.PaginationDataResponse, error)
	GetProductsByName(pageNumber int, name string, ctx context.Context) (response.PaginationDataResponse, error)
	GetProductsByDescription(pageNumber int, description string, ctx context.Context) (response.PaginationDataResponse, error)
	GetProductsByStatus(pageNumber int, status *bool, ctx context.Context) (response.PaginationDataResponse, error)
	GetProductsByKeyword(pageNumber int, keyword string, ctx context.Context) (response.PaginationDataResponse, error)
	GetProductById(id string, ctx context.Context) (*entity.Product, error)
	CreateProduct(req request.CreateProductRequest, ctx context.Context) error
	UpdateProduct(req request.UpdateProductRequest, ctx context.Context) error
	RemoveProduct(id string, ctx context.Context) error
	ActivateProduct(id string, ctx context.Context) error
}
