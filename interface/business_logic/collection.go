package businesslogic

import (
	"context"
	"sonit_server/model/dto/request"
	"sonit_server/model/entity"
)

type ICollectionService interface {
	GetAllCollections(ctx context.Context) (*[]entity.Collection, error)
	GetCollectionsByName(name string, ctx context.Context) (*[]entity.Collection, error)
	GetCollectionsByStatus(status *bool, ctx context.Context) (*[]entity.Collection, error)
	GetCollectionById(id string, ctx context.Context) (*entity.Collection, error)
	CreateCollection(req request.CreateCollectionRequest, ctx context.Context) error
	UpdateCollection(req request.UpdateCollectionRequest, ctx context.Context) error
	RemoveCollection(id string, ctx context.Context) error
	ActivateCollection(id string, ctx context.Context) error
}
