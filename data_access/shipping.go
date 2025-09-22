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

type shippingRepo struct {
	db     *sql.DB
	logger *log.Logger
}

func InitializeShippingRepo(db *sql.DB, logger *log.Logger) data_access.IShippingRepo {
	return &shippingRepo{
		db:     db,
		logger: logger,
	}
}

// CreateShipping implements dataaccess.IShippingRepo.
func (s *shippingRepo) CreateShipping(ship entity.Shipping, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetShippingTable()) + "CreateShipping - "
	var query string = "INSERT INTO " + entity.GetShippingTable() + " (id, delivery_code, shipping_unit, shipping_detail, delivered_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)"

	if _, err := s.db.Exec(query, ship.OrderId, ship.DeliveryCode, ship.ShippingUnit, ship.ShippingDetail, ship.DeliveredAt, ship.CreatedAt, ship.UpdatedAt); err != nil {
		s.logger.Println(errLogMsg + err.Error())
		return errors.New(noti.INTERNALL_ERR_MSG)
	}

	return nil
}

// GetShipping implements dataaccess.IShippingRepo.
func (s *shippingRepo) GetShipping(id string, ctx context.Context) (*entity.Shipping, int, error) {
	panic("unimplemented")
}

// GetUserShipping implements dataaccess.IShippingRepo.
func (s *shippingRepo) GetShippings(req request.GetShippingsRequest, ctx context.Context) (*[]entity.Shipping, int, error) {
	panic("unimplemented")
}

// UpdateShipping implements dataaccess.IShippingRepo.
func (s *shippingRepo) UpdateShipping(ship entity.Shipping, ctx context.Context) error {
	panic("unimplemented")
}

// RemoveShipping implements dataaccess.IShippingRepo.
func (s *shippingRepo) RemoveShipping(id string, ctx context.Context) error {
	var errLogMsg string = fmt.Sprintf(noti.REPO_ERR_MSG, entity.GetShippingTable()) + "RemoveShipping - "
	var query string = "DELETE FROM " + entity.GetShippingTable() + " WHERE id = $1"
	var INTERNALL_ERR_MSG error = errors.New(noti.INTERNALL_ERR_MSG)

	res, err := s.db.Exec(query, id)
	if err != nil {
		s.logger.Println(errLogMsg, err.Error())
		return INTERNALL_ERR_MSG
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		s.logger.Println(errLogMsg + err.Error())
		return INTERNALL_ERR_MSG
	}

	if rowsAffected == 0 {
		return errors.New(fmt.Sprintf(noti.UNDEFINED_OBJECT_WARN_MSG, entity.GetShippingTable()))
	}

	return nil
}
