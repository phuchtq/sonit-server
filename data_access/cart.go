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
)

type cartRepo struct {
	db     *sql.DB
	logger *log.Logger
}

func InitializeCartRepo(db *sql.DB, logger *log.Logger) data_access.ICartRepo {
	return &cartRepo{
		db:     db,
		logger: logger,
	}
}

const (
	items_in_cart_limit int = 10
)

// CreateCart implements repo.ICartRepo.
func (c *cartRepo) CreateCart(cart entity.Cart, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCartTable()) + "CreateCart - "
	var query string = "INSERT INTO " + entity.GetCartTable() + " (id, items, expired_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"

	if _, err := c.db.Exec(query, cart.UserId, cart.Items, cart.ExpiredAt, cart.CreatedAt, cart.UpdatedAt); err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetCartById implements repo.ICartRepo.
func (c *cartRepo) GetCart(id string, ctx context.Context) (*entity.Cart, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCartTable()) + "GetCart - "
	var query string = "SELECT * FROM " + entity.GetCartTable() + " WHERE id = $1"

	var res entity.Cart
	if err := c.db.QueryRow(query, id).Scan(&res.UserId, &res.Items, &res.ExpiredAt, &res.CreatedAt, &res.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, items_in_cart_limit, nil
		}

		c.logger.Println(errLogMsg + err.Error())
		return nil, items_in_cart_limit, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, items_in_cart_limit, nil
}

// RemoveCart implements repo.ICartRepo.
func (c *cartRepo) RemoveCart(id string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCartTable()) + "RemoveCart - "
	var query string = "DELETE FROM " + entity.GetCartTable() + " WHERE id = $1"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := c.db.Exec(query, id)
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
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetCartTable()))
	}

	return nil
}

// UpdateCart implements repo.ICartRepo.
func (c *cartRepo) UpdateCart(cart entity.Cart, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCartTable()) + "UpdateCart - "
	var query string = "UPDATE " + entity.GetCartTable() + " SET items = $1, expired_at = $2, updated_at = $3 WHERE id = $4"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := c.db.Exec(query, cart.Items, cart.ExpiredAt, cart.UpdatedAt, cart.UserId)
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
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetCartTable()))
	}

	return nil
}
