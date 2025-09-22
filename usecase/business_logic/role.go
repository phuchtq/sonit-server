package businesslogic

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sonit_server/constant/noti"
	repo "sonit_server/data_access" // Role data access ~~ Role repository
	"sonit_server/data_access/db"
	db_server "sonit_server/data_access/db_server"
	business_logic "sonit_server/interface/business_logic"
	data_access "sonit_server/interface/data_access" // Role repo interface
	"sonit_server/model/dto/request"
	entity "sonit_server/model/entity"
	"sonit_server/utils"
	"strings"
	"time"
)

type roleService struct {
	roleRepo data_access.IRoleRepo
	logger   *log.Logger
}

func GenerateRoleService() (business_logic.IRoleService, error) {
	var logger = utils.GetLogConfig()

	cnn, err := db.ConnectDB(logger, db_server.InitializePostgreSQL())

	if err != nil {
		return nil, err
	}

	role_cnn = cnn

	return InitializeRoleService(cnn, logger), nil
}

func InitializeRoleService(db *sql.DB, logger *log.Logger) business_logic.IRoleService {
	return &roleService{
		roleRepo: repo.InitializeRoleRepo(db, logger),
		logger:   logger,
	}
}

var role_cnn *sql.DB

const (
	id_type          string = "ID"
	name_type        string = "NAME"
	email_type       string = "EMAIL"
	description_type string = "DESCRIPTION"
)

// ActivateRole implements business_logic.IRoleService.
func (r *roleService) ActivateRole(id string, ctx context.Context) error {
	defer closeCnn(role_cnn)
	return r.roleRepo.ActivateRole(id, ctx)
}

// CreateRole implements business_logic.IRoleService.
func (r *roleService) CreateRole(name string, ctx context.Context) error {
	defer closeCnn(role_cnn)
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New(noti.FIELD_EMPTY_WARN_MSG)
	}
	//---------------------------------------
	if isEntityExist(r.roleRepo, name, name_type, ctx) {
		return errors.New(noti.ITEM_EXISTED_WARN_MSG)
	}
	//---------------------------------------
	var curTime time.Time = time.Now()
	return r.roleRepo.CreateRole(entity.Role{
		RoleId:       utils.GenerateId(),
		RoleName:     name,
		ActiveStatus: true,
		CreatedAt:    curTime,
		UpdatedAt:    curTime,
	}, ctx)
}

// GetAllRoles implements business_logic.IRoleService.
func (r *roleService) GetAllRoles(ctx context.Context) (*[]entity.Role, error) {
	defer closeCnn(role_cnn)
	return r.roleRepo.GetAllRoles(ctx)
}

// GetRoleById implements business_logic.IRoleService.
func (r *roleService) GetRoleById(id string, ctx context.Context) (*entity.Role, error) {
	defer closeCnn(role_cnn)

	if id == "" {
		return nil, errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}
	//---------------------------------------
	res, err := r.roleRepo.GetRoleById(id, ctx)
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, errors.New(entity.GetRoleTable() + noti.UNDEFINED_OBJECT_WARN_MSG)
	}
	//---------------------------------------
	return res, nil
}

// GetRolesByName implements business_logic.IRoleService.
func (r *roleService) GetRolesByName(name string, ctx context.Context) (*[]entity.Role, error) {
	defer closeCnn(role_cnn)

	if name == "" {
		return r.roleRepo.GetAllRoles(ctx)
	}

	return r.roleRepo.GetRolesByName(utils.ToNormalizedString(name), ctx)
}

// GetRolesByStatus implements business_logic.IRoleService.
func (r *roleService) GetRolesByStatus(status *bool, ctx context.Context) (*[]entity.Role, error) {
	defer closeCnn(role_cnn)

	if status == nil {
		return r.roleRepo.GetAllRoles(ctx)
	}

	return r.roleRepo.GetRolesByStatus(*status, ctx)
}

// RemoveRole implements business_logic.IRoleService.
func (r *roleService) RemoveRole(id string, ctx context.Context) error {
	defer closeCnn(role_cnn)

	if id == "" {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return r.roleRepo.RemoveRole(id, ctx)
}

// UpdateRole implements business_logic.IRoleService.
func (r *roleService) UpdateRole(req request.UpdateRoleRequest, ctx context.Context) error {
	defer closeCnn(role_cnn)

	role, err := r.roleRepo.GetRoleById(req.RoleId, ctx)
	if err != nil {
		return err
	}

	if req.RoleName == "" {
		return nil
	}

	role.RoleName = strings.TrimSpace(role.RoleName)
	role.UpdatedAt = time.Now()
	return r.roleRepo.UpdateRole(*role, ctx)
}
