package dataaccess

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sonit_server/constant/noti"
	data_access "sonit_server/interface/data_access"
	"sonit_server/model/dto/request"
	entity "sonit_server/model/entity"
)

type userSecurityRepo struct {
	db     *sql.DB
	logger *log.Logger
}

func InitializeUserSecurityRepo(db *sql.DB, logger *log.Logger) data_access.IUserSecurityRepo {
	return &userSecurityRepo{
		db:     db,
		logger: logger,
	}
}

// CreateUserSecurity implements repo.IUserSecurityRepo.
func (u *userSecurityRepo) CreateUserSecurity(usc entity.UserSecurity, ctx context.Context) error {
	var query string = "INSERT INTO " + entity.GetUserSecurityTable() + " (id, access_token, refresh_token, action_token, fail_access, last_fail) VALUES ($1, $2, $3, $4, $5, $6)"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetUserSecurityTable()) + "CreateUserSecurity - "

	if _, err := u.db.Exec(query, usc.UserId, usc.AccessToken, usc.RefreshToken, usc.ActionToken, usc.FailAccess, usc.LastFail); err != nil {
		u.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// Login implements repo.IUserSecurityRepo.
func (u *userSecurityRepo) Login(req request.LoginSecurityRequest, ctx context.Context) error {
	var query string = "UPDATE " + entity.GetUserSecurityTable() + " SET access_token = $1, refresh_token = $2 WHERE id = $3"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetUserSecurityTable()) + "Login - "
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := u.db.Exec(query, req.AccessToken, req.RefreshToken, req.UserId)
	if err != nil {
		u.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSG
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		u.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSG
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetUserSecurityTable()))
	}

	return nil
}

// Logout implements repo.IUserSecurityRepo.
func (u *userSecurityRepo) Logout(id string, ctx context.Context) error {
	var query string = "UPDATE " + entity.GetUserSecurityTable() +
		" SET access_token = NULL, refresh_token = NULL" +
		" WHERE id = $1"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetUserSecurityTable()) + "Logout - "
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)
	//defer u.db.Close()

	res, err := u.db.Exec(query, id)
	if err != nil {
		u.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSG
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		u.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSG
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetUserSecurityTable()))
	}

	return nil
}

// EditUserSecurity implements repo.IUserSecurityRepo.
func (u *userSecurityRepo) EditUserSecurity(usc entity.UserSecurity, ctx context.Context) error {
	var query string = "UPDATE " + entity.GetUserSecurityTable() + " SET access_token = $1, refresh_token = $2, action_token = $3, fail_access = $4, last_fail = $5 WHERE id = $6"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetUserSecurityTable()) + "EditUserSecurity - "
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := u.db.Exec(query, usc.AccessToken, usc.RefreshToken, usc.ActionToken, usc.FailAccess, &usc.LastFail, usc.UserId)
	if err != nil {
		u.logger.Println(errLogMsg, err.Error())
		return INTERNALL_ERR_MSG
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		u.logger.Println(errLogMsg, err.Error())
		return INTERNALL_ERR_MSG
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetUserSecurityTable()))
	}

	return nil
}

// GetUserSecurity implements repo.IUserSecurityRepo.
func (u *userSecurityRepo) GetUserSecurity(id string, ctx context.Context) (*entity.UserSecurity, error) {
	var query string = "SELECT * FROM " + entity.GetUserSecurityTable() + " WHERE id = $1"

	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetUserSecurityTable()) + "GetUserSecurity - "

	var usc entity.UserSecurity
	if err := u.db.QueryRow(query, id).Scan(
		&usc.UserId, &usc.AccessToken, &usc.RefreshToken,
		&usc.ActionToken, &usc.FailAccess, &usc.LastFail); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		u.logger.Println(errLogMsg, err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}
	return &usc, nil
}
