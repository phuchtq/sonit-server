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

type orderRepo struct {
	db     *sql.DB
	logger *log.Logger
}

func InitializeOrderRepo(db *sql.DB, logger *log.Logger) data_access.IOrderRepo {
	return &orderRepo{
		db:     db,
		logger: logger,
	}
}

const (
	items_in_order_limit int = 10
)

// CreateOrder implements dataaccess.IOrderRepo.
func (o *orderRepo) CreateOrder(order entity.Order, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetOrderTable()) + "CreateOrder - "
	var query string = "INSERT INTO " + entity.GetOrderTable() + " (id, user_id, items, total_amount, currency, status, note, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)"

	if _, err := o.db.Exec(query, order.OrderId, order.UserId, order.Items, order.TotalAmount, order.Currency, order.Status, order.Note, order.CreatedAt, order.UpdatedAt); err != nil {
		o.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetOrder implements dataaccess.IOrderRepo.
func (o *orderRepo) GetOrder(id string, ctx context.Context) (*entity.Order, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetOrderTable()) + "GetOrder - "
	var query string = "SELECT * FROM " + entity.GetOrderTable() + " WHERE id = $1"

	var res entity.Order
	if err := o.db.QueryRow(query, id).Scan(&res.OrderId, &res.UserId, &res.Items, &res.TotalAmount, &res.Currency, &res.Status, &res.Note, &res.CreatedAt, &res.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		o.logger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetOrders implements dataaccess.IOrderRepo.
func (o *orderRepo) GetOrders(req request.GetOrdersRequest, ctx context.Context) (*[]entity.Order, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetOrderTable()) + "GetOrders - "
	var queryCondition string = " WHERE "
	var isHavePreviousCond bool = false

	if req.UserId != "" {
		queryCondition += fmt.Sprintf("user_id = '%s'", req.UserId)
		isHavePreviousCond = true
	}
	if req.Request.Keyword != "" {
		if isHavePreviousCond {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(note) LIKE LOWER('%%%%%s%%%%')", req.Request.Keyword)
	}
	if req.Status != "" {
		if isHavePreviousCond {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(status) = LOWER('%%%%%s%%%%')", req.Status)
	}

	if queryCondition == " WHERE " {
		queryCondition = ""
	}

	var orderCondition string = generateOrderCondition(req.Request.FilterProp, req.Request.Order)
	var query string = generateRetrieveQuery(entity.GetOrderTable(), queryCondition+orderCondition, product_record_limit, req.Request.PageNumber, false)
	o.logger.Println("Query: ", query)

	rows, err := o.db.Query(query)
	if err != nil {
		o.logger.Println(errLogMsg + err.Error())
		return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var res []entity.Order
	for rows.Next() {
		var x entity.Order
		if err := rows.Scan(
			&x.OrderId, &x.UserId, &x.Items,
			&x.TotalAmount, &x.Currency, &x.Status, &x.Note,
			&x.CreatedAt, &x.UpdatedAt); err != nil {

			o.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	o.db.QueryRow(generateRetrieveQuery(entity.GetOrderTable(), queryCondition, product_record_limit, req.Request.PageNumber, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, product_record_limit), nil
}

// UpdateOrder implements dataaccess.IOrderRepo.
func (o *orderRepo) UpdateOrder(order entity.Order, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetOrderTable()) + "UpdateOrder - "
	var query string = "UPDATE " + entity.GetOrderTable() + " SET currency = $1, status = $2, note = $3, updated_at = $4 WHERE id = $5"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := o.db.Exec(query, order.Currency, order.Status, order.Note, order.UpdatedAt, order.OrderId)
	if err != nil {
		o.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSG
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		o.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSG
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetRoleTable()))
	}

	return nil
}
