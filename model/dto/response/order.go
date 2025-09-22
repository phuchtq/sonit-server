package response

import "time"

type ViewOrderResponse struct {
	OrderId     string     `json:"order_id"`
	UserId      string     `json:"user_id"`
	Items       []CartItem `json:"items"`
	TotalAmount float64    `json:"total_amount"`
	Currency    string     `json:"currency"`
	Status      string     `json:"status"`
	Note        string     `json:"note"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type ViewOrderCreatedSuccessReponse struct {
	OrderId     string     `json:"order_id"`
	Items       []CartItem `json:"items"`
	Status      string     `json:"status"`
	TotalAmount float64    `json:"total_amount"`
}
