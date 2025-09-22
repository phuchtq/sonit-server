package dataaccess

import (
	"context"
	entity "sonit_server/model/entity"
)

type IRoleRepo interface {
	GetAllRoles(ctx context.Context) (*[]entity.Role, error)
	GetRolesByName(name string, ctx context.Context) (*[]entity.Role, error)
	GetRolesByStatus(status bool, ctx context.Context) (*[]entity.Role, error)
	GetRoleById(id string, ctx context.Context) (*entity.Role, error)
	CreateRole(role entity.Role, ctx context.Context) error
	RemoveRole(id string, ctx context.Context) error
	UpdateRole(role entity.Role, ctx context.Context) error
	ActivateRole(id string, ctx context.Context) error
}
