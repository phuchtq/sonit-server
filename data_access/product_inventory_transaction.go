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

type productInventoryTransactionRepo struct {
	db     *sql.DB
	logger *log.Logger
}

func InitializeProductInventoryTransactionRepo(db *sql.DB, logger *log.Logger) data_access.IProductInventoryTransactionRepo {
	return &productInventoryTransactionRepo{
		db:     db,
		logger: logger,
	}
}

// CreateProductInventoryTransaction implements dataaccess.IProductInventoryTransactionRepo.
func (p *productInventoryTransactionRepo) CreateProductInventoryTransaction(tx entity.ProductInventoryTransaction, ctx context.Context) error {
	var query string = "INSERT INTO " + entity.GetProductInventoryTransactionTable() +
		" (id, product_id, amount, action, " +
		"note, date, created_at, updated at) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8)"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductInventoryTransactionTable()) + "CreateProductInventoryTransaction - "

	if _, err := p.db.Exec(query, tx.TransactionId, tx.ProductId, tx.Amount, tx.Action,
		tx.Note, tx.Date, tx.CreatedAt, tx.UpdatedAt); err != nil {

		p.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

const (
	product_inventory_transaction_records_limit int = 10
)

// GetAllProductInventoryTransactions implements dataaccess.IProductInventoryTransactionRepo.
func (p *productInventoryTransactionRepo) GetAllProductInventoryTransactions(req request.GetProductInventoryTrasactionsRequest, ctx context.Context) (*[]entity.ProductInventoryTransaction, int, error) {
	panic("unimplemented")
}

func (p *productInventoryTransactionRepo) GetInventoryTransactionsByProduct(req request.GetProductInventoryTrasactionsByProductRequest, ctx context.Context) (*[]entity.ProductInventoryTransaction, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductInventoryTransactionTable()) + "GetInventoryTransactionsByProduct - "
	var queryCondition string = fmt.Sprintf(` WHERE action = '%s' AND product_id = '%s'`, req.InventoryTransaction.Action, req.ProductId)
	var orderCondition string = generateOrderCondition(req.InventoryTransaction.Pagination.FilterProp, req.InventoryTransaction.Pagination.Order)
	var query string = generateRetrieveQuery(entity.GetProductInventoryTransactionTable(), queryCondition+orderCondition, product_record_limit, req.InventoryTransaction.Pagination.PageNumber, false)

	rows, err := p.db.Query(query)
	if err != nil {
		p.logger.Println(errLogMsg + err.Error())
		return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var res []entity.ProductInventoryTransaction
	for rows.Next() {
		var x entity.ProductInventoryTransaction
		if err := rows.Scan(
			&x.TransactionId, &x.ProductId, &x.Amount, &x.Action,
			&x.Note, &x.Date, &x.CreatedAt, &x.UpdatedAt); err != nil {

			p.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	p.db.QueryRow(generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_inventory_transaction_records_limit, req.InventoryTransaction.Pagination.PageNumber, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, product_record_limit), nil
}

// GetProductInventoryTransaction implements dataaccess.IProductInventoryTransactionRepo.
func (p *productInventoryTransactionRepo) GetProductInventoryTransaction(id string, ctx context.Context) (*entity.ProductInventoryTransaction, error) {
	var query string = "SELECT * FROM " + entity.GetProductInventoryTransactionTable() + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductInventoryTransactionTable()) + "GetProductInventoryTransaction - "

	var res entity.ProductInventoryTransaction
	if err := p.db.QueryRow(query, id).Scan(
		&res.TransactionId, &res.ProductId, &res.Amount, &res.Action,
		&res.Note, &res.Date, &res.CreatedAt, &res.UpdatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		p.logger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// UpdateProductInventoryTransaction implements dataaccess.IProductInventoryTransactionRepo.
func (p *productInventoryTransactionRepo) UpdateProductInventoryTransaction(tx entity.ProductInventoryTransaction, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductInventoryTransactionTable()) + "UpdateProductInventoryTransaction - "
	var query string = "UPDATE " + entity.GetProductInventoryTransactionTable() + " SET product_id = $1, amount = $2, action = $3, note = $4, date = $5, " +
		"updated_at = $6 WHERE id = $7"

	res, err := p.db.Exec(query, tx.ProductId, tx.Amount, tx.Action, tx.Note, tx.Date, tx.UpdatedAt, tx.TransactionId)

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
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetProductTable()))
	}

	return nil
}
