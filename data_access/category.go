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

type categoryRepo struct {
	db     *sql.DB
	logger *log.Logger
}

func InitializeCategoryRepo(db *sql.DB, logger *log.Logger) data_access.ICategoryRepo {
	return &categoryRepo{
		db:     db,
		logger: logger,
	}
}

// ActivateCategory implements repo.ICategoryRepo.
func (c *categoryRepo) ActivateCategory(id string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCategoryTable()) + "ActivateCategory - "
	var query string = "UPDATE " + entity.GetCategoryTable() + " SET active_status = true, updated_at = $1 WHERE id = $2"

	if _, err := c.db.Exec(query, fmt.Sprint(time.Now()), id); err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// CreateCategory implements repo.ICategoryRepo.
func (c *categoryRepo) CreateCategory(Category entity.Category, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCategoryTable()) + "CreateCategory - "
	var query string = "INSERT INTO " + entity.GetCategoryTable() + "(id, Category_name, description, active_status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"

	if _, err := c.db.Exec(query, Category.CategoryId, Category.CategoryName, Category.Description, Category.ActiveStatus, Category.CreatedAt, Category.UpdatedAt); err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetAllCategorys implements repo.ICategoryRepo.
func (c *categoryRepo) GetAllCategories(ctx context.Context) (*[]entity.Category, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCategoryTable()) + "GetAllCategories - "
	var query string = "SELECT * FROM " + entity.GetCategoryTable()
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := c.db.Query(query)
	if err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return nil, INTERNALL_ERR_MSG
	}

	var res []entity.Category
	for rows.Next() {
		var x entity.Category
		if err := rows.Scan(&x.CategoryId, &x.CategoryName, &x.Description, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {
			c.logger.Println(errLogMsg + err.Error())
			return nil, INTERNALL_ERR_MSG
		}

		res = append(res, x)
	}

	return &res, nil
}

// GetCategoryById implements repo.ICategoryRepo.
func (c *categoryRepo) GetCategoryById(id string, ctx context.Context) (*entity.Category, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCategoryTable()) + "GetCategoryById - "
	var query string = "SELECT * FROM " + entity.GetCategoryTable() + " WHERE id = $1"

	var res entity.Category
	if err := c.db.QueryRow(query, id).Scan(&res.CategoryId, &res.Description, &res.CategoryName, &res.ActiveStatus, &res.CreatedAt, &res.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		c.logger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetCategorysByName implements repo.ICategoryRepo.
func (c *categoryRepo) GetCategoriesByName(name string, ctx context.Context) (*[]entity.Category, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCategoryTable()) + "GetCategorysByName - "
	var query string = "SELECT * FROM " + entity.GetCategoryTable() + " WHERE LOWER(name) = LOWER('%$1%')"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := c.db.Query(query, name)
	if err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return nil, INTERNALL_ERR_MSG
	}

	var res []entity.Category
	for rows.Next() {
		var x entity.Category
		if err := rows.Scan(&x.CategoryId, &x.CategoryName, &x.Description, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {
			c.logger.Println(errLogMsg, err.Error())
			return nil, INTERNALL_ERR_MSG
		}

		res = append(res, x)
	}

	return &res, nil
}

// GetCategorysByStatus implements repo.ICategoryRepo.
func (c *categoryRepo) GetCategoriesByStatus(status bool, ctx context.Context) (*[]entity.Category, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCategoryTable()) + "GetCategorysByStatus - "
	var query string = "SELECT * FROM " + entity.GetCategoryTable() + " WHERE active_status = $1"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := c.db.Query(query, fmt.Sprint(status))
	if err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return nil, INTERNALL_ERR_MSG
	}

	var res []entity.Category
	for rows.Next() {
		var x entity.Category
		if err := rows.Scan(&x.CategoryId, &x.CategoryName, &x.Description, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {
			c.logger.Println(errLogMsg + err.Error())
			return nil, INTERNALL_ERR_MSG
		}

		res = append(res, x)
	}

	return &res, nil
}

// RemoveCategory implements repo.ICategoryRepo.
func (c *categoryRepo) RemoveCategory(id string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCategoryTable()) + "RemoveCategory - "
	var query string = "UPDATE " + entity.GetCategoryTable() + " SET active_status = false, updated_at = $1 WHERE id = $2"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := c.db.Exec(query, time.Now().String(), id)
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
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetCategoryTable()))
	}

	return nil
}

// UpdateCategory implements repo.ICategoryRepo.
func (c *categoryRepo) UpdateCategory(category entity.Category, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCategoryTable()) + "UpdateCategory - "
	var query string = "UPDATE " + entity.GetCategoryTable() + " SET category_name = $1, description = $2, active_status = $3, updated_at = $4 WHERE id = $5"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := c.db.Exec(query, category.CategoryName, category.Description, category.ActiveStatus, category.UpdatedAt, category.CategoryId)
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
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetCategoryTable()))
	}

	return nil
}
