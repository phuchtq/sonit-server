package dataaccess

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sonit_server/constant/noti"
	data_access "sonit_server/interface/data_access"
	entity "sonit_server/model/entity"
	"time"
)

type roleRepo struct {
	db     *sql.DB
	logger *log.Logger
}

func InitializeRoleRepo(db *sql.DB, logger *log.Logger) data_access.IRoleRepo {
	return &roleRepo{
		db:     db,
		logger: logger,
	}
}

// ActivateRole implements repo.IRoleRepo.
func (r *roleRepo) ActivateRole(id string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetRoleTable()) + "ActivateRole - "
	var query string = "UPDATE " + entity.GetRoleTable() + " SET active_status = true, updated_at = $1 WHERE id = $2"

	if _, err := r.db.Exec(query, fmt.Sprint(time.Now()), id); err != nil {
		r.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// CreateRole implements repo.IRoleRepo.
func (r *roleRepo) CreateRole(role entity.Role, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetRoleTable()) + "CreateRole - "
	var query string = "INSERT INTO " + entity.GetRoleTable() + " (id, name, active_status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"

	if _, err := r.db.Exec(query, role.RoleId, role.RoleName, role.ActiveStatus, role.CreatedAt, role.UpdatedAt); err != nil {
		r.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetAllRoles implements repo.IRoleRepo.
func (r *roleRepo) GetAllRoles(ctx context.Context) (*[]entity.Role, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetRoleTable()) + "GetAllRoles - "
	var query string = "SELECT * FROM " + entity.GetRoleTable()
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := r.db.Query(query)
	if err != nil {
		r.logger.Println(errLogMsg + err.Error())
		return nil, INTERNALL_ERR_MSG
	}

	var res []entity.Role
	for rows.Next() {
		var x entity.Role
		if err := rows.Scan(&x.RoleId, &x.RoleName, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {
			r.logger.Println(errLogMsg + err.Error())
			return nil, INTERNALL_ERR_MSG
		}

		res = append(res, x)
	}

	return &res, nil
}

// GetRoleById implements repo.IRoleRepo.
func (r *roleRepo) GetRoleById(id string, ctx context.Context) (*entity.Role, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetRoleTable()) + "GetRoleById - "
	var query string = "SELECT * FROM " + entity.GetRoleTable() + " WHERE id = $1"

	var res entity.Role
	if err := r.db.QueryRow(query, id).Scan(&res.RoleId, &res.RoleName, &res.ActiveStatus, &res.CreatedAt, &res.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		r.logger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetRolesByName implements repo.IRoleRepo.
func (r *roleRepo) GetRolesByName(name string, ctx context.Context) (*[]entity.Role, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetRoleTable()) + "GetRolesByName - "
	var query string = "SELECT * FROM " + entity.GetRoleTable() + " WHERE LOWER(name) = LOWER('%$1%')"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := r.db.Query(query, name)
	if err != nil {
		r.logger.Println(errLogMsg + err.Error())
		return nil, INTERNALL_ERR_MSG
	}

	var res []entity.Role
	for rows.Next() {
		var x entity.Role
		if err := rows.Scan(&x.RoleId, &x.RoleName, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {
			r.logger.Println(errLogMsg, err.Error())
			return nil, INTERNALL_ERR_MSG
		}

		res = append(res, x)
	}

	return &res, nil
}

// GetRolesByStatus implements repo.IRoleRepo.
func (r *roleRepo) GetRolesByStatus(status bool, ctx context.Context) (*[]entity.Role, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetRoleTable()) + "GetRolesByStatus - "
	var query string = "SELECT * FROM " + entity.GetRoleTable() + " WHERE active_status = $1"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := r.db.Query(query, fmt.Sprint(status))
	if err != nil {
		r.logger.Println(errLogMsg + err.Error())
		return nil, INTERNALL_ERR_MSG
	}

	var res []entity.Role
	for rows.Next() {
		var x entity.Role
		if err := rows.Scan(&x.RoleId, &x.RoleName, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {
			r.logger.Println(errLogMsg + err.Error())
			return nil, INTERNALL_ERR_MSG
		}

		res = append(res, x)
	}

	return &res, nil
}

// RemoveRole implements repo.IRoleRepo.
func (r *roleRepo) RemoveRole(id string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetRoleTable()) + "RemoveRole - "
	var query string = "UPDATE " + entity.GetRoleTable() + " SET active_status = false, updated_at = $1 WHERE id = $2"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := r.db.Exec(query, time.Now().String(), id)
	if err != nil {
		r.logger.Println(errLogMsg, err.Error())
		return INTERNALL_ERR_MSG
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSG
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetRoleTable()))
	}

	return nil
}

// UpdateRole implements repo.IRoleRepo.
func (r *roleRepo) UpdateRole(role entity.Role, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetRoleTable()) + "UpdateRole - "
	var query string = "UPDATE " + entity.GetRoleTable() + " SET name = $1, active_status = $2, updated_at = $3 WHERE id = $4"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := r.db.Exec(query, role.RoleName, role.ActiveStatus, role.UpdatedAt, role.RoleId)
	if err != nil {
		r.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSG
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSG
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetRoleTable()))
	}

	return nil
}
