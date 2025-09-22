package businesslogic

import (
	"context"
	"sonit_server/model/dto/request"
	"sonit_server/model/entity"
)

type IVoucherService interface {
	GetAllVouchers(ctx context.Context) (*[]entity.Voucher, error)
	GetAllValidVouchers(ctx context.Context) (*[]entity.Voucher, error)
	GetVoucherByID(id string, ctx context.Context) (*entity.Voucher, error)
	GetVoucherByCode(code string, ctx context.Context) (*entity.Voucher, error)
	CreateVoucher(req request.CreateVoucherRequest, ctx context.Context) error
	UpdateVoucher(req request.UpdateVoucherRequest, ctx context.Context) error
	RemoveVoucher(id string, ctx context.Context) error
}
