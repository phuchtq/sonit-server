package request

import "time"

type CreateVoucherRequest struct {
	Code               string    `json:"code"`
	Discount           float64   `json:"discount"` // e.g. 10.0 for 10%
	Amount             int64     `json:"amount"`
	Description        string    `json:"description"`
	AllowedCategoryIDs []string  `json:"allowed_category_ids"`
	AllowedProductIDs  []string  `json:"allowed_product_ids"`
	ExpiredAt          time.Time `json:"expires_at"`
}

type UpdateVoucherRequest struct {
	VoucherID          string    `json:"voucher_id" validate:"required"`
	Code               string    `json:"code"`
	Discount           float64   `json:"discount"` // e.g. 10.0 for 10%
	Amount             int64     `json:"amount"`
	Description        string    `json:"description"`
	ActiveStatus       bool      `json:"active_status"`
	AllowedCategoryIDs []string  `json:"allowed_category_ids"`
	AllowedProductIDs  []string  `json:"allowed_product_ids"`
	ExpiredAt          time.Time `json:"expires_at"`
}
