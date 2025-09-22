package dataaccess

import (
	"context"
	"sonit_server/model/entity"
)

type ICategoryRepo interface {
	GetAllCategories(ctx context.Context) (*[]entity.Category, error)
	GetCategoriesByName(name string, ctx context.Context) (*[]entity.Category, error)
	GetCategoriesByStatus(status bool, ctx context.Context) (*[]entity.Category, error)
	GetCategoryById(id string, ctx context.Context) (*entity.Category, error)
	CreateCategory(Category entity.Category, ctx context.Context) error
	UpdateCategory(Category entity.Category, ctx context.Context) error
	RemoveCategory(id string, ctx context.Context) error
	ActivateCategory(id string, ctx context.Context) error
}
