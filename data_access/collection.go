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

type collectionRepo struct {
	db     *sql.DB
	logger *log.Logger
}

func InitializeCollectionRepo(db *sql.DB, logger *log.Logger) data_access.ICollectionRepo {
	return &collectionRepo{
		db:     db,
		logger: logger,
	}
}

// ActivateCollection implements repo.ICollectionRepo.
func (c *collectionRepo) ActivateCollection(id string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCollectionTable()) + "ActivateCollection - "
	var query string = "UPDATE " + entity.GetCollectionTable() + " SET active_status = true, updated_at = $1 WHERE id = $2"

	if _, err := c.db.Exec(query, fmt.Sprint(time.Now()), id); err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// CreateCollection implements repo.ICollectionRepo.
func (c *collectionRepo) CreateCollection(collection entity.Collection, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCollectionTable()) + "CreateCollection - "
	var query string = "INSERT INTO " + entity.GetCollectionTable() + "(id, collection_name, description, active_status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)"

	if _, err := c.db.Exec(query, collection.CollectionId, collection.CollectionName, collection.Description, collection.ActiveStatus, collection.CreatedAt, collection.UpdatedAt); err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetAllCollections implements repo.ICollectionRepo.
func (c *collectionRepo) GetAllCollections(ctx context.Context) (*[]entity.Collection, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCollectionTable()) + "GetAllCollections - "
	var query string = "SELECT * FROM " + entity.GetCollectionTable()
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := c.db.Query(query)
	if err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return nil, INTERNALL_ERR_MSG
	}

	var res []entity.Collection
	for rows.Next() {
		var x entity.Collection
		if err := rows.Scan(&x.CollectionId, &x.CollectionName, &x.Description, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {
			c.logger.Println(errLogMsg + err.Error())
			return nil, INTERNALL_ERR_MSG
		}

		res = append(res, x)
	}

	return &res, nil
}

// GetCollectionById implements repo.ICollectionRepo.
func (c *collectionRepo) GetCollectionById(id string, ctx context.Context) (*entity.Collection, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCollectionTable()) + "GetCollectionById - "
	var query string = "SELECT * FROM " + entity.GetCollectionTable() + " WHERE id = $1"

	var res entity.Collection
	if err := c.db.QueryRow(query, id).Scan(&res.CollectionId, &res.CollectionName, &res.Description, &res.ActiveStatus, &res.CreatedAt, &res.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		c.logger.Println(errLogMsg + err.Error())
		return nil, errors.New(noti.INTERNALL_ERR_MSG)
	}

	return &res, nil
}

// GetCollectionsByName implements repo.ICollectionRepo.
func (c *collectionRepo) GetCollectionsByName(name string, ctx context.Context) (*[]entity.Collection, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCollectionTable()) + "GetCollectionsByName - "
	var query string = "SELECT * FROM " + entity.GetCollectionTable() + " WHERE LOWER(name) = LOWER('%$1%')"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := c.db.Query(query, name)
	if err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return nil, INTERNALL_ERR_MSG
	}

	var res []entity.Collection
	for rows.Next() {
		var x entity.Collection
		if err := rows.Scan(&x.CollectionId, &x.CollectionName, &x.Description, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {
			c.logger.Println(errLogMsg, err.Error())
			return nil, INTERNALL_ERR_MSG
		}

		res = append(res, x)
	}

	return &res, nil
}

// GetCollectionsByStatus implements repo.ICollectionRepo.
func (c *collectionRepo) GetCollectionsByStatus(status bool, ctx context.Context) (*[]entity.Collection, error) {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCollectionTable()) + "GetCollectionsByStatus - "
	var query string = "SELECT * FROM " + entity.GetCollectionTable() + " WHERE active_status = $1"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	rows, err := c.db.Query(query, fmt.Sprint(status))
	if err != nil {
		c.logger.Println(errLogMsg + err.Error())
		return nil, INTERNALL_ERR_MSG
	}

	var res []entity.Collection
	for rows.Next() {
		var x entity.Collection
		if err := rows.Scan(&x.CollectionId, &x.CollectionName, &x.Description, &x.ActiveStatus, &x.CreatedAt, &x.UpdatedAt); err != nil {
			c.logger.Println(errLogMsg + err.Error())
			return nil, INTERNALL_ERR_MSG
		}

		res = append(res, x)
	}

	return &res, nil
}

// RemoveCollection implements repo.ICollectionRepo.
func (c *collectionRepo) RemoveCollection(id string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCollectionTable()) + "RemoveCollection - "
	var query string = "UPDATE " + entity.GetCollectionTable() + " SET active_status = false, updated_at = $1 WHERE id = $2"
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
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetCollectionTable()))
	}

	return nil
}

// UpdateCollection implements repo.ICollectionRepo.
func (c *collectionRepo) UpdateCollection(collection entity.Collection, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetCollectionTable()) + "UpdateCollection - "
	var query string = "UPDATE " + entity.GetCollectionTable() + " SET collection_name = $1, description = $2, active_status = $3, updated_at = $4 WHERE id = $5"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := c.db.Exec(query, collection.CollectionName, collection.Description, collection.ActiveStatus, collection.UpdatedAt, collection.CollectionId)
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
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetCollectionTable()))
	}

	return nil
}
