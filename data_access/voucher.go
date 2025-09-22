package dataaccess

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"

	"log"
	"sonit_server/constant/noti"
	data_access "sonit_server/interface/data_access"
	entity "sonit_server/model/entity"
	"time"
)

type voucherRepo struct {
	db     *sql.DB
	logger *log.Logger
}

func InitializeVoucherRepo(db *sql.DB, logger *log.Logger) data_access.IVoucherRepo {
	return &voucherRepo{
		db:     db,
		logger: logger,
	}
}

// GetAllVouchers implements repo.IVoucherRepo.
func (c *voucherRepo) GetAllVouchers(ctx context.Context) (*[]entity.Voucher, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetVoucherTable()) + "GetAllVouchers - "
	var query string = "SELECT * FROM " + entity.GetVoucherTable()
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := c.db.Query(query)
	if err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return nil, INTERNALL_ERR_MSG
	}

	var res []entity.Voucher
	for rows.Next() {
		var x entity.Voucher
		if err := rows.Scan(&x.VoucherId, &x.Code, &x.Discount, &x.Amount, &x.Description, &x.ActiveStatus,
			pq.Array(&x.AllowedCategoryIDs), pq.Array(&x.AllowedProductIDs), &x.ExpiredAt, &x.CreatedAt, &x.UpdatedAt); err != nil {
			c.logger.Println(errLogMsg + err.Error())
			return nil, INTERNALL_ERR_MSG
		}

		res = append(res, x)
	}

	return &res, nil
}

// GetAllValidVouchers implements repo.IVoucherRepo.
func (c *voucherRepo) GetAllValidVouchers(ctx context.Context) (*[]entity.Voucher, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetVoucherTable()) + "GetAllValidVouchers - "
	var query string = "SELECT * FROM " + entity.GetVoucherTable() +
		" WHERE active_status = TRUE AND expired_at > NOW() AND amount > 0"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := c.db.Query(query)
	if err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return nil, INTERNALL_ERR_MSG
	}

	var res []entity.Voucher
	for rows.Next() {
		var x entity.Voucher
		if err := rows.Scan(&x.VoucherId, &x.Code, &x.Discount, &x.Amount, &x.Description, &x.ActiveStatus,
			pq.Array(&x.AllowedCategoryIDs), pq.Array(&x.AllowedProductIDs), &x.ExpiredAt, &x.CreatedAt, &x.UpdatedAt); err != nil {
			c.logger.Println(errLogMsg + err.Error())
			return nil, INTERNALL_ERR_MSG
		}

		res = append(res, x)
	}

	return &res, nil
}

// CreateVoucher implements repo.IVoucherRepo.
func (c *voucherRepo) CreateVoucher(voucher entity.Voucher, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetVoucherTable()) + "CreateVoucher - "
	var query string = "INSERT INTO " + entity.GetVoucherTable() + "(id, code, discount, amount, description, " +
		"active_status, allowed_category_ids, allowed_product_ids, expired_at, created_at, updated_at) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"

	if _, err := c.db.Exec(query, voucher.VoucherId, voucher.Code, voucher.Discount, voucher.Amount,
		voucher.Description, voucher.ActiveStatus, voucher.AllowedCategoryIDs, voucher.AllowedProductIDs,
		voucher.ExpiredAt, voucher.CreatedAt, voucher.UpdatedAt); err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetVoucherByID implements repo.IVoucherRepo.
func (c *voucherRepo) GetVoucherByID(id string, ctx context.Context) (*entity.Voucher, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetVoucherTable()) + "GetVoucherByID - "
	var query string = "SELECT * FROM " + entity.GetVoucherTable() + " WHERE id = $1"

	var res *entity.Voucher
	if err := c.db.QueryRow(query, id).Scan(&res.VoucherId, &res.Code, &res.Discount, &res.Amount,
		&res.Description, &res.ActiveStatus, pq.Array(&res.AllowedCategoryIDs), pq.Array(&res.AllowedProductIDs),
		&res.ExpiredAt, &res.CreatedAt, &res.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		c.logger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return res, nil
}

// GetVoucherByCode implements repo.IVoucherRepo.
func (c *voucherRepo) GetVoucherByCode(code string, ctx context.Context) (*entity.Voucher, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetVoucherTable()) + "GetVoucherByCode - "
	var query string = "SELECT * FROM " + entity.GetVoucherTable() + " WHERE LOWER(code) = LOWER($1)"

	var res *entity.Voucher
	if err := c.db.QueryRow(query, code).Scan(&res.VoucherId, &res.Code, &res.Discount, &res.Amount,
		&res.Description, &res.ActiveStatus, pq.Array(&res.AllowedCategoryIDs), pq.Array(&res.AllowedProductIDs),
		&res.ExpiredAt, &res.CreatedAt, &res.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		c.logger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return res, nil
}

// UpdateVoucher implements repo.IVoucherRepo.
func (c *voucherRepo) UpdateVoucher(voucher entity.Voucher, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetVoucherTable()) + "UpdateVoucher - "
	var query string = "UPDATE " + entity.GetVoucherTable() + " SET code = $1, discount = $2, amount = $3," +
		"description = $4, active_status = $5, allowed_category_ids = $6, allowed_product_ids = $7, " +
		"expired_at = $8, updated_at = $9 WHERE id = $10"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := c.db.Exec(query, voucher.Code, voucher.Discount, voucher.Amount, voucher.Description,
		voucher.ActiveStatus, pq.Array(voucher.AllowedCategoryIDs), pq.Array(voucher.AllowedProductIDs),
		voucher.ExpiredAt, time.Now(), voucher.VoucherId)
	if err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSG
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSG
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetVoucherTable()))
	}

	return nil
}

// RemoveVoucher implements repo.IVoucherRepo.
func (c *voucherRepo) RemoveVoucher(id string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetVoucherTable()) + "RemoveCategory - "
	var query string = "UPDATE " + entity.GetVoucherTable() + " SET active_status = FALSE, updated_at = $1 WHERE id = $2"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := c.db.Exec(query, time.Now(), id)
	if err != nil {
		c.logger.Println(errLogMsg, err.Error())
		return INTERNALL_ERR_MSG
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSG
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetVoucherTable()))
	}

	return nil
}
