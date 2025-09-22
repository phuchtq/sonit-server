package dataaccess

import (
	"context"
	entity "sonit_server/model/entity"
)

type IUserRepo interface {
	GetAllUsers(pageNumber int, ctx context.Context) (*[]entity.User, int, error)
	GetUsersByRole(pageNumber int, id string, ctx context.Context) (*[]entity.User, int, error)
	GetUsersByStatus(pageNumber int, status bool, ctx context.Context) (*[]entity.User, int, error)
	GetUser(id string, ctx context.Context) (*entity.User, error)
	GetUserByEmail(email string, ctx context.Context) (*entity.User, error)
	GetUsersByKeyword(pageNumber int, keyword string, ctx context.Context) (*[]entity.User, int, error)
	CreateUser(user entity.User, ctx context.Context) error
	UpdateUser(user entity.User, ctx context.Context) error
	ChangeUserStatus(id string, status bool, ctx context.Context) error
	IsVipCodeExist(code string, ctx context.Context) (bool, error)
}
