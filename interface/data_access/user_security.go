package dataaccess

import (
	"context"
	"sonit_server/model/dto/request"
	entity "sonit_server/model/entity"
)

type IUserSecurityRepo interface {
	GetUserSecurity(id string, ctx context.Context) (*entity.UserSecurity, error)
	CreateUserSecurity(usc entity.UserSecurity, ctx context.Context) error
	EditUserSecurity(usc entity.UserSecurity, ctx context.Context) error
	Login(req request.LoginSecurityRequest, ctx context.Context) error
	Logout(id string, ctx context.Context) error
}
