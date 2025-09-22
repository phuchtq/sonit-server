package request

type GetPaymentsRequest struct {
	Request SearchPaginatioRequest `json:"request"`
	Status  string                 `json:"status" form:"status"`
	Method  string                 `json:"method" form:"method"`
	UserId  string                 `json:"user_id" form:"user_id"`
}

type UpdatePaymentRequest struct {
	PaymentId string `json:"payment_id" validate:"required"`
	Currency  string `json:"currency"`
	Status    string `json:"status"`
	Method    string `json:"method"`
}

type CreatePaymentDirectRequest struct {
	UserId      string       `json:"user_id" validate:"required"`
	Note        string       `json:"note"`
	Product     ItemSelected `json:"product" validate:"required"`
	Address     string       `jsosn:"address" validate:"required"`
	PhoneNumber string       `json:"phone_number" validate:"required"`
}

type CreatePaymentThroughCartRequest struct {
	UserId      string         `json:"user_id" validate:"required"`
	Note        string         `json:"note"`
	Items       []ItemSelected `json:"items" validate:"required"`
	Address     string         `jsosn:"address" validate:"required"`
	PhoneNumber string         `json:"phone_number" validate:"required"`
}
