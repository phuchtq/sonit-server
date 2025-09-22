package entity

import "time"

type Cart struct {
	UserId    string    `json:"user_id"`
	Items     string    `json:"items"`
	ExpiredAt time.Time `json:"expired_at" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetCartTable() string {
	return "carts"
}
