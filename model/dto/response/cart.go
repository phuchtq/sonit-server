package response

import "time"

type CartItem struct {
	ProductId string  `json:"product_id"`
	Name      string  `json:"name"`
	ImageUrl  string  `json:"image_url"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Currency  string  `json:"currency"`
}

type ViewGeneralCartUIResponse struct {
	UserId   string `json:"user_id"`
	ImageUrl string `json:"image_url"`
	Quantity int    `json:"quantity"`
}

type ViewCartResponse struct {
	UserId    string     `json:"user_id"`
	Items     []CartItem `json:"items"`
	ExpiredAt time.Time  `json:"expired_at"`
}
