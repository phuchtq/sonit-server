package dataaccess

import (
	"context"
	"sonit_server/model/dto/request"
	"sonit_server/model/entity"
)

type IProductRepo interface {
	GetAllProducts(pageNumber int, ctx context.Context) (*[]entity.Product, int, error)
	GetProductsCustomerUI(req request.GetProductsCustomerUI, ctx context.Context) (*[]entity.Product, int, error)
	GetProductsByCategory(pageNumber int, id string, ctx context.Context) (*[]entity.Product, int, error)
	GetProductsByCollection(pageNumber int, id string, ctx context.Context) (*[]entity.Product, int, error)
	GetProductsByPriceInterval(pageNumber int, maxPrice, minPrice int64, ctx context.Context) (*[]entity.Product, int, error)
	GetProductsByName(pageNumber int, name string, ctx context.Context) (*[]entity.Product, int, error)
	GetProductsByDescription(pageNumber int, description string, ctx context.Context) (*[]entity.Product, int, error)
	GetProductsByStatus(pageNumber int, status bool, ctx context.Context) (*[]entity.Product, int, error)
	GetProductsByKeyword(pageNumber int, keyword string, ctx context.Context) (*[]entity.Product, int, error)
	GetProductById(id string, ctx context.Context) (*entity.Product, error)
	CreateProduct(product entity.Product, ctx context.Context) error
	UpdateProduct(product entity.Product, ctx context.Context) error
	RemoveProduct(id string, ctx context.Context) error
	ActivateProduct(id string, ctx context.Context) error
}
