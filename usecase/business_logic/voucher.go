package businesslogic

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"sonit_server/constant/noti"
	repo "sonit_server/data_access"
	"sonit_server/data_access/db"
	db_server "sonit_server/data_access/db_server"
	business_logic "sonit_server/interface/business_logic"
	data_access "sonit_server/interface/data_access"
	"sonit_server/model/dto/request"
	"sonit_server/model/entity"
	"sonit_server/utils"
	"time"
)

type voucherService struct {
	voucherRepo data_access.IVoucherRepo
	logger      *log.Logger
}

func GenerateVoucherService() (business_logic.IVoucherService, error) {
	var logger = utils.GetLogConfig()

	cnn, err := db.ConnectDB(logger, db_server.InitializePostgreSQL())

	if err != nil {
		return nil, err
	}

	voucher_cnn = cnn

	return InitializeVoucherService(cnn, logger), nil
}

func InitializeVoucherService(db *sql.DB, logger *log.Logger) business_logic.IVoucherService {
	return &voucherService{
		voucherRepo: repo.InitializeVoucherRepo(db, logger),
		logger:      logger,
	}
}

var voucher_cnn *sql.DB

func (v voucherService) GetAllVouchers(ctx context.Context) (*[]entity.Voucher, error) {
	defer closeCnn(voucher_cnn)
	return v.voucherRepo.GetAllVouchers(ctx)
}

func (v voucherService) GetAllValidVouchers(ctx context.Context) (*[]entity.Voucher, error) {
	defer closeCnn(voucher_cnn)
	return v.voucherRepo.GetAllValidVouchers(ctx)
}

func (v voucherService) GetVoucherByID(id string, ctx context.Context) (*entity.Voucher, error) {
	defer closeCnn(voucher_cnn)

	if id == "" {
		return nil, errors.New(noti.FIELD_EMPTY_WARN_MSG)
	}

	return v.voucherRepo.GetVoucherByID(utils.ToNormalizedString(id), ctx)
}

func (v voucherService) GetVoucherByCode(code string, ctx context.Context) (*entity.Voucher, error) {
	defer closeCnn(voucher_cnn)

	if code == "" {
		return nil, errors.New(noti.FIELD_EMPTY_WARN_MSG)
	}

	return v.voucherRepo.GetVoucherByCode(utils.ToNormalizedString(code), ctx)
}

func (v voucherService) CreateVoucher(req request.CreateVoucherRequest, ctx context.Context) error {
	defer closeCnn(voucher_cnn)

	// Collection name existed
	if isEntityExist(v.voucherRepo, req.Code, name_type, ctx) {
		return errors.New(noti.ITEM_EXISTED_WARN_MSG)
	}
	//---------------------------------------
	var curTime time.Time = time.Now()
	return v.voucherRepo.CreateVoucher(entity.Voucher{
		VoucherId:          utils.GenerateId(),
		Code:               req.Code,
		Discount:           req.Discount,
		Amount:             req.Amount,
		Description:        req.Description,
		ActiveStatus:       true,
		AllowedCategoryIDs: req.AllowedCategoryIDs,
		AllowedProductIDs:  req.AllowedProductIDs,
		ExpiredAt:          req.ExpiredAt,
		CreatedAt:          curTime,
		UpdatedAt:          curTime,
	}, ctx)
}

func (v voucherService) UpdateVoucher(req request.UpdateVoucherRequest, ctx context.Context) error {
	defer closeCnn(voucher_cnn)

	voucher, err := v.voucherRepo.GetVoucherByID(req.VoucherID, ctx)
	if err != nil {
		return err
	}

	if req.Code != "" && !isEntityExist(v.voucherRepo, req.Code, name_type, ctx) {
		voucher.Code = req.Code
	}

	// Update Description: only if req.Description is non-empty.
	if req.Description != "" {
		voucher.Description = req.Description
	}

	// Update ExpiredAt: only if req.ExpiredAt is a non-zero time.
	if !req.ExpiredAt.IsZero() {
		voucher.ExpiredAt = req.ExpiredAt
	}

	// Update other fields directly. If these fields are omitted in JSON,
	voucher.Discount = req.Discount
	voucher.Amount = req.Amount
	voucher.ActiveStatus = req.ActiveStatus
	voucher.AllowedCategoryIDs = req.AllowedCategoryIDs
	voucher.AllowedProductIDs = req.AllowedProductIDs

	voucher.UpdatedAt = time.Now()

	return v.voucherRepo.UpdateVoucher(*voucher, ctx)

}

func (v voucherService) RemoveVoucher(id string, ctx context.Context) error {
	defer closeCnn(voucher_cnn)

	if id == "" {
		return errors.New(noti.GENERIC_ERROR_WARN_MSG)
	}

	return v.voucherRepo.RemoveVoucher(id, ctx)
}
