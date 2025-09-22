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

type productRepo struct {
	db     *sql.DB
	logger *log.Logger
}

func InitializeProductRepo(db *sql.DB, logger *log.Logger) data_access.IProductRepo {
	return &productRepo{
		db:     db,
		logger: logger,
	}
}

const (
	product_record_limit int = 10
)

// GetAllProducts implements repo.IProductRepo.
func (p *productRepo) GetAllProducts(pageNumber int, ctx context.Context) (*[]entity.Product, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductTable()) + "GetAllProducts - "
	var query string = generateRetrieveQuery(entity.GetProductTable(), "", product_record_limit, pageNumber, false)

	rows, err := p.db.Query(query)
	if err != nil {
		p.logger.Println(errLogMsg + err.Error())
		return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var res []entity.Product
	for rows.Next() {
		var x entity.Product
		if err := rows.Scan(
			&x.ProductId, &x.CategoryId, &x.CollectionId, &x.ProductName, &x.Description,
			&x.Image, &x.Size, &x.Color, &x.Price,
			&x.Currency, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {

			p.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	p.db.QueryRow(generateRetrieveQuery(entity.GetProductTable(), "", product_record_limit, pageNumber, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, product_record_limit), nil
}

// GetProductsCustomerUI implements dataaccess.IProductRepo.
func (p *productRepo) GetProductsCustomerUI(req request.GetProductsCustomerUI, ctx context.Context) (*[]entity.Product, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductTable()) + "GetProductsCustomerUI - "
	var queryCondition string = " WHERE active_status = true"
	if req.CategoryId != "" {
		queryCondition += fmt.Sprintf(" AND category_id = '%s'", req.CategoryId)
	}
	if req.CollectionId != "" {
		queryCondition += fmt.Sprintf(" AND collection_id = '%s'", req.CollectionId)
	}
	if req.Pagination.Keyword != "" {
		queryCondition += fmt.Sprintf(" AND (LOWER(description) LIKE LOWER('%%%%%s%%%%') OR LOWER(name) LIKE LOWER('%%%%%s%%%%'))", req.Pagination.Keyword, req.Pagination.Keyword)
	}

	var orderCondition string = generateOrderCondition(req.Pagination.FilterProp, req.Pagination.Order)
	var query string = generateRetrieveQuery(entity.GetProductTable(), queryCondition+orderCondition, product_record_limit, req.Pagination.PageNumber, false)

	p.logger.Println("Query: ", query)

	rows, err := p.db.Query(query)
	if err != nil {
		p.logger.Println(errLogMsg + err.Error())
		return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var res []entity.Product
	for rows.Next() {
		var x entity.Product
		if err := rows.Scan(
			&x.ProductId, &x.CategoryId, &x.CollectionId, &x.ProductName, &x.Description,
			&x.Image, &x.Size, &x.Color, &x.Price,
			&x.Currency, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {

			p.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	p.db.QueryRow(generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_record_limit, req.Pagination.PageNumber, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, product_record_limit), nil
}

// GetProductsByStatus implements repo.IProductRepo.
func (p *productRepo) GetProductsByStatus(pageNumber int, status bool, ctx context.Context) (*[]entity.Product, int, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductTable()) + "GetProductsByStatus - "
	var queryCondition string = fmt.Sprintf(" WHERE is_active = %t", status)
	var query string = generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_record_limit, pageNumber, false)

	rows, err := p.db.Query(query)
	if err != nil {
		p.logger.Println(errLogMsg + err.Error())
		return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var res []entity.Product
	for rows.Next() {
		var x entity.Product
		if err := rows.Scan(
			&x.ProductId, &x.CategoryId, &x.CollectionId, &x.ProductName, &x.Description,
			&x.Image, &x.Size, &x.Color, &x.Price,
			&x.Currency, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {

			p.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	p.db.QueryRow(generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_record_limit, 1, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, product_record_limit), nil
}

// GetProductsByCategory implements repo.IProductRepo.
func (p *productRepo) GetProductsByCategory(pageNumber int, id string, ctx context.Context) (*[]entity.Product, int, error) {
	var queryCondition string = fmt.Sprintf(" WHERE category_id = %s", id)
	var query string = generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_record_limit, pageNumber, false)
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductTable()) + "GetProductsByCategory - "

	rows, err := p.db.Query(query)
	if err != nil {
		p.logger.Println(errLogMsg + err.Error())
		return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var res []entity.Product
	for rows.Next() {
		var x entity.Product
		if err := rows.Scan(
			&x.ProductId, &x.CategoryId, &x.CollectionId, &x.ProductName, &x.Description,
			&x.Image, &x.Size, &x.Color, &x.Price,
			&x.Currency, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {

			p.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	p.db.QueryRow(generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_record_limit, pageNumber, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, product_record_limit), nil
}

// GetProductsByCollection implements repo.IProductRepo.
func (p *productRepo) GetProductsByCollection(pageNumber int, id string, ctx context.Context) (*[]entity.Product, int, error) {
	var queryCondition string = fmt.Sprintf(" WHERE collection_id = '%s'", id)
	var query string = generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_record_limit, pageNumber, false)
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductTable()) + "GetProductsByCollection - "

	rows, err := p.db.Query(query)
	if err != nil {
		p.logger.Println(errLogMsg + err.Error())
		return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var res []entity.Product
	for rows.Next() {
		var x entity.Product
		if err := rows.Scan(
			&x.ProductId, &x.CategoryId, &x.CollectionId, &x.ProductName, &x.Description,
			&x.Image, &x.Size, &x.Color, &x.Price,
			&x.Currency, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {

			p.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	p.db.QueryRow(generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_record_limit, pageNumber, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, product_record_limit), nil
}

// GetProductsByPriceInterval implements repo.IProductRepo.
func (p *productRepo) GetProductsByPriceInterval(pageNumber int, maxPrice, minPrice int64, ctx context.Context) (*[]entity.Product, int, error) {
	var queryCondition string = fmt.Sprintf(" WHERE price >= %d AND price <= %d", minPrice, maxPrice)
	var query string = generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_record_limit, pageNumber, false)
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductTable()) + "GetProductsByPriceInterval - "

	rows, err := p.db.Query(query)
	if err != nil {
		p.logger.Println(errLogMsg + err.Error())
		return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var res []entity.Product
	for rows.Next() {
		var x entity.Product
		if err := rows.Scan(
			&x.ProductId, &x.CategoryId, &x.CollectionId, &x.ProductName, &x.Description,
			&x.Image, &x.Size, &x.Color, &x.Price,
			&x.Currency, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {

			p.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	p.db.QueryRow(generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_record_limit, pageNumber, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, product_record_limit), nil
}

// GetProductsByName implements repo.IProductRepo.
func (p *productRepo) GetProductsByName(pageNumber int, name string, ctx context.Context) (*[]entity.Product, int, error) {
	var queryCondition string = fmt.Sprintf(" WHERE LOWER(name) LIKE LOWER('%s')", name)
	var query string = generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_record_limit, pageNumber, false)
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductTable()) + "GetProductsByName - "

	rows, err := p.db.Query(query)
	if err != nil {
		p.logger.Println(errLogMsg + err.Error())
		return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var res []entity.Product
	for rows.Next() {
		var x entity.Product
		if err := rows.Scan(
			&x.ProductId, &x.CategoryId, &x.CollectionId, &x.ProductName, &x.Description,
			&x.Image, &x.Size, &x.Color, &x.Price,
			&x.Currency, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {

			p.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	p.db.QueryRow(generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_record_limit, pageNumber, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, product_record_limit), nil
}

// GetProductsByDescription implements repo.IProductRepo.
func (p *productRepo) GetProductsByDescription(pageNumber int, description string, ctx context.Context) (*[]entity.Product, int, error) {
	var queryCondition string = fmt.Sprintf(" WHERE LOWER(name) LIKE LOWER('%s')", description)
	var query string = generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_record_limit, pageNumber, false)
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductTable()) + "GetProductsByDescription - "

	rows, err := p.db.Query(query)
	if err != nil {
		p.logger.Println(errLogMsg + err.Error())
		return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
	}

	var res []entity.Product
	for rows.Next() {
		var x entity.Product
		if err := rows.Scan(
			&x.ProductId, &x.CategoryId, &x.CollectionId, &x.ProductName, &x.Description,
			&x.Image, &x.Size, &x.Color, &x.Price,
			&x.Currency, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {

			p.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	p.db.QueryRow(generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_record_limit, pageNumber, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, product_record_limit), nil
}

// UpdateProduct implements repo.IProductRepo.
func (p *productRepo) UpdateProduct(product entity.Product, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductTable()) + "UpdateProduct - "
	var query string = "UPDATE " + entity.GetProductTable() + " SET category_id = $1,  collection_id = $2, name = $3, description = $4, " +
		"image = $5, size = $6, color = $7, price = $8, currency = $9, " +
		"active_status = $10, updated_at = $11 " +
		"WHERE id = $12"

	res, err := p.db.Exec(query, product.CategoryId, product.CollectionId, product.ProductName, product.Description,
		product.Image, product.Size, product.Color, product.Price, product.Currency,
		product.ActiveStatus, product.UpdatedAt, product.ProductId)

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

// CreateProduct implements repo.IProductRepo.
func (p *productRepo) CreateProduct(product entity.Product, ctx context.Context) error {
	var query string = "INSERT INTO " + entity.GetProductTable() +
		"(id, category_id, collection_id, name, description, " +
		"image, size, color, price, currency, " +
		"active_status, created_at, updated at) " +
		"values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductTable()) + "CreateProduct - "

	if _, err := p.db.Exec(query, product.ProductId, product.CategoryId, product.CollectionId, product.ProductName, product.Description,
		product.Image, product.Size, product.Color, product.Price, product.Currency,
		product.ActiveStatus, product.CreatedAt, product.UpdatedAt); err != nil {

		p.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetProduct implements repo.IProductRepo.
func (p *productRepo) GetProductById(id string, ctx context.Context) (*entity.Product, error) {
	var query string = "SELECT * FROM " + entity.GetProductTable() + " WHERE id = $1"
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductTable()) + "GetProductById - "

	p.logger.Println("Query: ", query)

	var res entity.Product
	if err := p.db.QueryRow(query, id).Scan(
		&res.ProductId, &res.CategoryId, &res.CollectionId, &res.ProductName, &res.Description,
		&res.Image, &res.Size, &res.Color, &res.Price, &res.Currency,
		&res.ActiveStatus, &res.CreatedAt, &res.UpdatedAt); err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		p.logger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetProductsByKeyword implements repo.IProductRepo.
func (p *productRepo) GetProductsByKeyword(pageNumber int, keyword string, ctx context.Context) (*[]entity.Product, int, error) {
	var queryCondition string = fmt.Sprintf(` WHERE LOWER(name) LIKE LOWER('%%' || %s || '%%') OR LOWER(description) LIKE LOWER('%%' || %s || '%%') OR LOWER(size) LIKE LOWER('%%' || %s || '%%')`, keyword, keyword, keyword)
	var query string = generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_record_limit, pageNumber, false)
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductTable()) + "GetProductsByKeyword - "
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := p.db.Query(query, keyword, keyword)
	if err != nil {
		p.logger.Println(errLogMsg + err.Error())
		return nil, 0, INTERNALL_ERR_MSG
	}

	var res []entity.Product
	for rows.Next() {
		var x entity.Product
		if err := rows.Scan(
			&x.ProductId, &x.CategoryId, &x.CollectionId, &x.ProductName, &x.Description,
			&x.Image, &x.Size, &x.Color, &x.Price,
			&x.Currency, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {

			p.logger.Println(errLogMsg + err.Error())
			return nil, 0, errors.New(noti.INTERNALL_ERR_MSG)
		}

		res = append(res, x)
	}

	// Track total records in table
	var totalRecords int
	p.db.QueryRow(generateRetrieveQuery(entity.GetProductTable(), queryCondition, product_record_limit, pageNumber, true)).Scan(&totalRecords)

	return &res, caculateTotalPages(totalRecords, product_record_limit), nil
}

// ActivateProduct implements dataaccess.IProductRepo.
func (p *productRepo) ActivateProduct(id string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductTable()) + "ActivateProduct - "
	var query string = "UPDATE " + entity.GetProductTable() + " SET active_status = true WHERE id = $1"

	res, err := p.db.Exec(query, id)

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

// RemoveProduct implements dataaccess.IProductRepo.
func (p *productRepo) RemoveProduct(id string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetProductTable()) + "RemoveProduct - "
	var query string = "UPDATE " + entity.GetProductTable() + " SET active_status = false WHERE id = $1"

	res, err := p.db.Exec(query, id)

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
