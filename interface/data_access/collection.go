package dataaccess

import (
	"context"
	"sonit_server/model/entity"
)

type ICollectionRepo interface {
	GetAllCollections(ctx context.Context) (*[]entity.Collection, error)
	GetCollectionsByName(name string, ctx context.Context) (*[]entity.Collection, error)
	GetCollectionsByStatus(status bool, ctx context.Context) (*[]entity.Collection, error)
	GetCollectionById(id string, ctx context.Context) (*entity.Collection, error)
	CreateCollection(collection entity.Collection, ctx context.Context) error
	UpdateCollection(collection entity.Collection, ctx context.Context) error
	RemoveCollection(id string, ctx context.Context) error
	ActivateCollection(id string, ctx context.Context) error
}
