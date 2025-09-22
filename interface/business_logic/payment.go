package businesslogic

import (
	"context"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
	entity "sonit_server/model/entity"
)

type IPaymentService interface {
	GetPayments(req request.GetPaymentsRequest, ctx context.Context) (response.PaginationDataResponse, error)
	GetPaymentById(id string, ctx context.Context) (*entity.Payment, error)
	CreatePaymentDirect(req request.CreatePaymentDirectRequest, ctx context.Context) (response.UrlAPIResponse, error)
	CreatePaymentThroughCart(req request.CreatePaymentThroughCartRequest, ctx context.Context) (response.UrlAPIResponse, error)
	UpdatePayment(req request.UpdatePaymentRequest, ctx context.Context) error
	// Callback function
	CallbackPaymentSuccess(id string, ctx context.Context) (string, error)
	CallbackPaymentCancel(id string, ctx context.Context) (string, error)
}
