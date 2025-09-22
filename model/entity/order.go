package entity

import "time"

type Order struct {
	OrderId     string    `json:"order_id"`
	UserId      string    `json:"user_id"`
	Items       string    `json:"items"`
	TotalAmount float64   `json:"total_amount"`
	Currency    string    `json:"currency"`
	Status      string    `json:"status"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GetOrderTable() string {
	return "orders"
}
