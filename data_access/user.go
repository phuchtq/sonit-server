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
	"sync"
	"time"
)

type userRepo struct {
	db     *sql.DB
	logger *log.Logger
}

func InitializeUserRepo(db *sql.DB, logger *log.Logger) data_access.IUserRepo {
	return &userRepo{
		db:     db,
		logger: logger,
	}
}

const (
	user_record_limit int = 10
)

// IsVipCodeExist implements dataaccess.IUserRepo.
func (u *userRepo) IsVipCodeExist(code string, ctx context.Context) (bool, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetUserTable()) + "IsVipCodeExist - "
	var record int

	if err := u.db.QueryRow(generateCountTotalRecordsQuery(entity.GetUserTable(), " WHERE vip_code = "+code)).Scan(&record); err != nil {
		u.logger.Println(errLogMsg + err.Error())
		return false, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return record > 0, nil
}

// ChangeUserStatus implements repo.IUserRepo.
func (u *userRepo) ChangeUserStatus(id string, status bool, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetUserTable()) + "ChangeUserStatus - "

	var lastFailValueQuery string = "NULL"
	if status {
		lastFailValueQuery = fmt.Sprint(time.Date(1900, time.January, 1, 0, 0, 0, 0, time.UTC))
	}

	var userQuery string = "UPDATE " + entity.GetUserTable() + " SET is_active = $1 AND updated_at = $2 WHERE user_id = $3"
	var securityQuery string = "UPDATE " + entity.GetUserSecurityTable() + " SET access_token = NULL, access_expiration = NULL," +
		" refresh_token = NULL, refresh_expiration = NULL, action_token = NULL, action_expiration = NULL, fail_access = 0, last_fail = " +
		lastFailValueQuery + " WHERE id = $1"
	//defer u.db.Close()

	var errChan chan error = make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		//defer wg.Done()
		if _, err := u.db.Exec(userQuery, fmt.Sprint(status), time.Now().GoString(), id); err != nil {
			u.logger.Println(errLogMsg + err.Error())
			errChan <- err
		}
	}()

	go func() {
		//defer wg.Done()
		if _, err := u.db.Exec(securityQuery, id); err != nil {
			u.logger.Println(errLogMsg + err.Error())
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)
	for err := range errChan {
		if err != nil {
			return errors.New(noti.INTERNALL_ERR_MSG)
		}
	}

	return nil
}

// GetAllUsers implements repo.IUserRepo.
func (u *userRepo) GetAllUsers(pageNumber int, ctx context.Context) (*[]entity.User, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetUserTable()) + "GetAllUsers - "
	var query string = generateRetrieveQuery(entity.GetUserTable(), "", user_record_limit, pageNumber, false)

	rows, err := u.db.Query(query)
	if err != nil {
		u.logger.Println(errLogMsg + err.Error())
		return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var res []entity.User
	for rows.Next() {
		var x entity.User
		var isHaveToResetPw sql.NullBool
		var vipCode sql.NullString

		if err := rows.Scan(
			&x.UserId, &x.RoleId, &x.FullName, &x.Email, &x.Password,
			&x.ProfileAvatar, &x.Gender, &x.IsVip, &vipCode,
			&x.IsActive, &x.IsActivated, &isHaveToResetPw, &x.CreatedAt, &x.UpdatedAt); err != nil {

			u.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		if isHaveToResetPw.Valid {
			x.IsHaveToResetPw = &isHaveToResetPw.Bool
		}

		if vipCode.Valid {
			x.VipCode = &vipCode.String
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	u.db.QueryRow(generateRetrieveQuery(entity.GetUserTable(), "", user_record_limit, pageNumber, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, user_record_limit), nil
}

// GetUsersByStatus implements repo.IUserRepo.
func (u *userRepo) GetUsersByStatus(pageNumber int, status bool, ctx context.Context) (*[]entity.User, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetUserTable()) + "GetUsersByStatus - "
	var queryCondition string = fmt.Sprintf(" WHERE is_active = %t", status)
	var query string = generateRetrieveQuery(entity.GetUserTable(), queryCondition, user_record_limit, pageNumber, false)

	rows, err := u.db.Query(query)
	if err != nil {
		u.logger.Println(errLogMsg + err.Error())
		return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var res []entity.User
	for rows.Next() {
		var x entity.User
		var isHaveToResetPw sql.NullBool
		var vipCode sql.NullString

		if err := rows.Scan(
			&x.UserId, &x.RoleId, &x.FullName, &x.Email, &x.Password,
			&x.ProfileAvatar, &x.Gender, &x.IsVip, &vipCode,
			&x.IsActive, &x.IsActivated, &isHaveToResetPw, &x.CreatedAt, &x.UpdatedAt); err != nil {

			u.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		if isHaveToResetPw.Valid {
			x.IsHaveToResetPw = &isHaveToResetPw.Bool
		}

		if vipCode.Valid {
			x.VipCode = &vipCode.String
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	u.db.QueryRow(generateRetrieveQuery(entity.GetUserTable(), queryCondition, user_record_limit, -1, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, user_record_limit), nil
}

// GetUserByEmail implements repo.IUserRepo.
func (u *userRepo) GetUserByEmail(email string, ctx context.Context) (*entity.User, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetUserTable()) + "GetUserByEmail - "
	var query string = "SELECT * from " + entity.GetUserTable() + " WHERE LOWER(email) = LOWER($1)"

	var res entity.User
	var isHaveToResetPw sql.NullBool
	var vipCode sql.NullString

	if err := u.db.QueryRow(query, email).Scan(
		&res.UserId, &res.RoleId, &res.FullName, &res.Email, &res.Password,
		&res.ProfileAvatar, &res.Gender, &res.IsVip, &vipCode,
		&res.IsActive, &res.IsActivated, &isHaveToResetPw, &res.CreatedAt, &res.UpdatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		u.logger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	// Transfer nullable values to the result struct
	if isHaveToResetPw.Valid {
		res.IsHaveToResetPw = &isHaveToResetPw.Bool
	}

	if vipCode.Valid {
		res.VipCode = &vipCode.String
	}

	return &res, nil
}

// GetUsersByRole implements repo.IUserRepo.
func (u *userRepo) GetUsersByRole(pageNumber int, id string, ctx context.Context) (*[]entity.User, int, error) {
	var queryCondition string = fmt.Sprintf(" WHERE role_id = %s", id)
	var query string = generateRetrieveQuery(entity.GetUserTable(), queryCondition, user_record_limit, pageNumber, false)
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetUserTable()) + "GetUsersByRole - "

	rows, err := u.db.Query(query, id)
	if err != nil {
		u.logger.Println(errLogMsg + err.Error())
		return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var res []entity.User
	for rows.Next() {
		var x entity.User
		var isHaveToResetPw sql.NullBool
		var vipCode sql.NullString

		if err := rows.Scan(
			&x.UserId, &x.RoleId, &x.FullName, &x.Email, &x.Password,
			&x.ProfileAvatar, &x.Gender, &x.IsVip, &vipCode,
			&x.IsActive, &x.IsActivated, &isHaveToResetPw, &x.CreatedAt, &x.UpdatedAt); err != nil {

			u.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		if isHaveToResetPw.Valid {
			x.IsHaveToResetPw = &isHaveToResetPw.Bool
		}

		if vipCode.Valid {
			x.VipCode = &vipCode.String
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	u.db.QueryRow(generateRetrieveQuery(entity.GetUserTable(), queryCondition, user_record_limit, pageNumber, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, user_record_limit), nil
}

// UpdateUser implements repo.IUserRepo.
func (u *userRepo) UpdateUser(user entity.User, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetUserTable()) + "UpdateUser - "
	var query string = "UPDATE " + entity.GetUserTable() + " SET email = $1, role_id = $2, full_name = $3, password = $4, " +
		"profile_avatar = $5, gender = $6, is_vip = $7, vip_code = $8, " +
		"is_active = $9, is_activated = $10, is_have_to_reset_password = $11, " +
		"updated_at = $12 WHERE id = $13"

	res, err := u.db.Exec(query, user.Email, user.RoleId, user.FullName, user.Password,
		user.ProfileAvatar, user.Gender, user.IsVip, user.VipCode,
		user.IsActive, user.IsActivated, user.IsHaveToResetPw, user.UpdatedAt, user.UserId)

	var INTERNALL_ERR_MSGMsg error = errors.New(noti.INTERNALL_ERR_MSG)

	if err != nil {
		u.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSGMsg
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		u.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSGMsg
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetUserTable()))
	}

	return nil
}

// CreateUser implements repo.IUserRepo.
func (u *userRepo) CreateUser(user entity.User, ctx context.Context) error {
	var query string = "INSERT INTO " + entity.GetUserTable() +
		"(id, role_id, full_name, email, password, profile_avatar, " +
		"gender, is_vip, vip_code, " +
		"is_active, is_activated, is_have_to_reset_password, " +
		"created_at, updated_at) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetUserTable()) + "CreateUser - "

	if _, err := u.db.Exec(query, user.UserId, user.RoleId, user.FullName, user.Email, user.Password,
		user.ProfileAvatar, user.Gender, user.IsVip, user.VipCode,
		user.IsActive, user.IsActivated, user.IsHaveToResetPw, user.CreatedAt, user.UpdatedAt); err != nil {

		u.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetUser implements repo.IUserRepo.
func (u *userRepo) GetUser(id string, ctx context.Context) (*entity.User, error) {
	var query string = "SELECT * FROM " + entity.GetUserTable() + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetUserTable()) + "GetUser - "

	var res entity.User
	var isHaveToResetPw sql.NullBool
	var vipCode sql.NullString

	if err := u.db.QueryRow(query, id).Scan(
		&res.UserId, &res.RoleId, &res.FullName, &res.Email, &res.Password,
		&res.ProfileAvatar, &res.Gender, &res.IsVip, &vipCode,
		&res.IsActive, &res.IsActivated, &isHaveToResetPw, &res.CreatedAt, &res.UpdatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		u.logger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	// Transfer nullable values to the result struct
	if isHaveToResetPw.Valid {
		res.IsHaveToResetPw = &isHaveToResetPw.Bool
	}

	if vipCode.Valid {
		res.VipCode = &vipCode.String
	}

	return &res, nil
}

// GetUsersByKeyword implements repo.IUserRepo.
func (u *userRepo) GetUsersByKeyword(pageNumber int, keyword string, ctx context.Context) (*[]entity.User, int, error) {
	var queryCondition string = fmt.Sprintf(` WHERE LOWER(full_name) LIKE LOWER('%%' || %s || '%%') OR LOWER(email) LIKE LOWER('%%' || %s || '%%')`, keyword, keyword)
	var query string = generateRetrieveQuery(entity.GetUserTable(), queryCondition, user_record_limit, pageNumber, false)
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetUserTable()) + "GetUsersByKeyword - "
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := u.db.Query(query, keyword, keyword)
	if err != nil {
		u.logger.Println(errLogMsg + err.Error())
		return nil, 0, INTERNALL_ERR_MSG
	}

	var res []entity.User
	for rows.Next() {
		var x entity.User
		var isHaveToResetPw sql.NullBool
		var vipCode sql.NullString

		if err := rows.Scan(
			&x.UserId, &x.RoleId, &x.FullName, &x.Email, &x.Password,
			&x.ProfileAvatar, &x.Gender, &x.IsVip, &vipCode,
			&x.IsActive, &x.IsActivated, &isHaveToResetPw, &x.CreatedAt, &x.UpdatedAt); err != nil {

			u.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		if isHaveToResetPw.Valid {
			x.IsHaveToResetPw = &isHaveToResetPw.Bool
		}

		if vipCode.Valid {
			x.VipCode = &vipCode.String
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	u.db.QueryRow(generateRetrieveQuery(entity.GetUserTable(), queryCondition, user_record_limit, pageNumber, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, user_record_limit), nil
}
