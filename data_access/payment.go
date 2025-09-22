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
	"sonit_server/model/entity"
)

type paymentRepo struct {
	db     *sql.DB
	logger *log.Logger
}

func InitializePaymentRepo(db *sql.DB, logger *log.Logger) data_access.IPaymentRepo {
	return &paymentRepo{
		db:     db,
		logger: logger,
	}
}

const (
	payment_limit_records int = 10
)

// CreatePayment implements dataaccess.IPaymentRepo.
func (p *paymentRepo) CreatePayment(payment entity.Payment, ctx context.Context) error {
	var query string = "INSERT INTO " + entity.GetPaymentTable() +
		"(id, user_id, order_id, transaction_id, " +
		"amount, currency, status, method, " +
		"created_at, updated_at) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetPaymentTable()) + "CreatePayment - "

	if _, err := p.db.Exec(query, payment.PaymentId, payment.UserId, payment.OrderId, payment.TransactionId,
		payment.Amount, payment.Currency, payment.Status, payment.Method,
		payment.CreatedAt, payment.UpdatedAt); err != nil {

		p.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetAllPayments implements dataaccess.IPaymentRepo.
func (p *paymentRepo) GetPayments(req request.GetPaymentsRequest, ctx context.Context) (*[]entity.Payment, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetPaymentTable()) + "GetPayments - "
	// var queryCondition string = fmt.Sprintf(
	// 	" WHERE status = '%s' AND method = '%s' AND user_id = '%s'",
	// 	req.Status,
	// 	req.Method,
	// 	req.UserId,
	// )
	var queryCondition string = " WHERE "
	var isHavePreviousCond bool = false

	if req.Status != "" {
		queryCondition += fmt.Sprintf("LOWER(status) = LOWER('%%%%%s%%%%')", req.Status)
		isHavePreviousCond = true
	}

	if req.UserId != "" {
		if isHavePreviousCond {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("user_id = '%s'", req.UserId)
		isHavePreviousCond = true
	}
	if req.Method != "" {
		if isHavePreviousCond {
			queryCondition += " AND "
		}

		queryCondition += fmt.Sprintf("LOWER(method) = LOWER('%%%%%s%%%%')", req.Method)
	}

	if queryCondition == " WHERE " {
		queryCondition = ""
	}

	var orderCondition string = generateOrderCondition(req.Request.FilterProp, req.Request.Order)
	var query string = generateRetrieveQuery(entity.GetPaymentTable(), queryCondition+orderCondition, payment_limit_records, req.Request.PageNumber, false)
	p.logger.Println("Query: ", query)

	rows, err := p.db.Query(query)
	if err != nil {
		p.logger.Println(errLogMsg + err.Error())
		return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var res []entity.Payment
	for rows.Next() {
		var x entity.Payment
		if err := rows.Scan(
			&x.PaymentId, &x.UserId, &x.OrderId, &x.TransactionId,
			&x.Amount, &x.Currency, &x.Status, &x.Method,
			&x.CreatedAt, &x.UpdatedAt); err != nil {

			p.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	p.db.QueryRow(generateRetrieveQuery(entity.GetPaymentTable(), queryCondition, payment_limit_records, req.Request.PageNumber, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, payment_limit_records), nil
}

// GetPaymentById implements dataaccess.IPaymentRepo.
func (p *paymentRepo) GetPaymentById(id string, ctx context.Context) (*entity.Payment, error) {
	var query string = "SELECT * FROM " + entity.GetPaymentTable() + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetPaymentTable()) + "GetPaymentById - "

	var res entity.Payment
	if err := p.db.QueryRow(query, id).Scan(
		&res.PaymentId, &res.UserId, &res.OrderId, &res.TransactionId,
		&res.Amount, &res.Currency, &res.Status, &res.Method,
		&res.CreatedAt, &res.UpdatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		p.logger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// UpdatePayment implements dataaccess.IPaymentRepo.
func (p *paymentRepo) UpdatePayment(payment entity.Payment, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetPaymentTable()) + "UpdatePayment - "
	var query string = "UPDATE " + entity.GetPaymentTable() + " SET status = $1,  method = $2, updated_at = $3 WHERE id = $4"

	res, err := p.db.Exec(query, payment.Status, payment.Method, payment.UpdatedAt, payment.PaymentId)

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
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetPaymentTable()))
	}

	return nil
}
