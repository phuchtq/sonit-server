package entity

import "time"

type Payment struct {
	PaymentId     string    `json:"payment_id"`
	UserId        string    `json:"user_id"`
	OrderId       string    `json:"order_id"`
	TransactionId string    `json:"transaction_id"`
	Amount        float64   `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	Method        string    `json:"method"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func GetPaymentTable() string {
	return "payments"
}
