package businesslogic

import (
	"context"
	"sonit_server/model/dto/request"
	entity "sonit_server/model/entity"
)

type IRoleService interface {
	GetAllRoles(ctx context.Context) (*[]entity.Role, error)
	GetRolesByName(name string, ctx context.Context) (*[]entity.Role, error)
	GetRolesByStatus(status *bool, ctx context.Context) (*[]entity.Role, error)
	GetRoleById(id string, ctx context.Context) (*entity.Role, error)
	CreateRole(name string, ctx context.Context) error
	UpdateRole(req request.UpdateRoleRequest, ctx context.Context) error
	RemoveRole(id string, ctx context.Context) error
	ActivateRole(id string, ctx context.Context) error
}
