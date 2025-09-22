package dataaccess

import (
	"context"
	"sonit_server/model/dto/request"
	"sonit_server/model/entity"
)

type IPaymentRepo interface {
	GetPayments(req request.GetPaymentsRequest, ctx context.Context) (*[]entity.Payment, int, error)
	GetPaymentById(id string, ctx context.Context) (*entity.Payment, error)
	CreatePayment(payment entity.Payment, ctx context.Context) error
	UpdatePayment(payment entity.Payment, ctx context.Context) error
}
