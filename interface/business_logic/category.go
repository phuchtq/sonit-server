package businesslogic

import (
	"context"
	"sonit_server/model/dto/request"
	entity "sonit_server/model/entity"
)

type ICategoryService interface {
	GetAllCategories(ctx context.Context) (*[]entity.Category, error)
	GetCategoriesByName(name string, ctx context.Context) (*[]entity.Category, error)
	GetCategoriesByStatus(status *bool, ctx context.Context) (*[]entity.Category, error)
	GetCategoryById(id string, ctx context.Context) (*entity.Category, error)
	CreateCategory(req request.CreateCategoryRequest, ctx context.Context) error
	UpdateCategory(req request.UpdateCategoryRequest, ctx context.Context) error
	RemoveCategory(id string, ctx context.Context) error
	ActivateCategory(id string, ctx context.Context) error
}
