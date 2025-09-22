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

type productInventoryRepo struct {
	db     *sql.DB
	logger *log.Logger
}

func InitializeProductInventoryRepo(db *sql.DB, logger *log.Logger) data_access.IProductInventoryRepo {
	return &productInventoryRepo{
		db:     db,
		logger: logger,
	}
}

// GetProductInventory implements dataaccess.IProductInventtory.
func (p *productInventoryRepo) GetProductInventory(id string, ctx context.Context) (*entity.ProductInventory, error) {
	var query string = "SELECT * FROM " + entity.GetProductInventoryTable() + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductTable()) + "GetProductInventory - "

	var res entity.ProductInventory
	if err := p.db.QueryRow(query, id).Scan(&res.ProductId, &res.CurrentQuantity); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		p.logger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// UpdateProductInventory implements dataaccess.IProductInventtory.
func (p *productInventoryRepo) UpdateProductInventory(inventory entity.ProductInventory, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductInventoryTable()) + "UpdateProductInventory - "
	var query string = "UPDATE " + entity.GetProductInventoryTable() + " SET current_quantity = $1 WHERE id = $2"

	res, err := p.db.Exec(query, inventory.CurrentQuantity, inventory.ProductId)

	var INTERNALL_ERR_MSGMsg error = errors.New(noti.INTERNALL_ERR_MSG)

	if err != nil {
		p.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSGMsg
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		p.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSGMsg
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetProductInventoryTable()))
	}

	return nil
}

// CreateProductInventory implements dataaccess.IProductInventtory.
func (p *productInventoryRepo) CreateProductInventory(inventory entity.ProductInventory, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductInventoryTable()) + "CreateProductInventory - "
	var query string = "INSERT INTO " + entity.GetRoleTable() + " (id, current_quantity) VALUES ($1, $2)"

	if _, err := p.db.Exec(query, inventory.ProductId, inventory.CurrentQuantity); err != nil {
		p.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}
