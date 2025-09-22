package dataaccess

import (
	"golang.org/x/net/context"
	"sonit_server/model/entity"
)

type IVoucherRepo interface {
	GetAllVouchers(ctx context.Context) (*[]entity.Voucher, error)
	GetAllValidVouchers(ctx context.Context) (*[]entity.Voucher, error)
	GetVoucherByID(id string, ctx context.Context) (*entity.Voucher, error)
	GetVoucherByCode(code string, ctx context.Context) (*entity.Voucher, error)
	CreateVoucher(voucher entity.Voucher, ctx context.Context) error
	UpdateVoucher(voucher entity.Voucher, ctx context.Context) error
	RemoveVoucher(id string, ctx context.Context) error
}
