package entity

import "time"

type Voucher struct {
	VoucherId          string    `json:"voucher_id"`
	Code               string    `json:"code"`
	Discount           float64   `json:"discount" validate:"required"` // e.g. 10.0 for 10%
	Amount             int64     `json:"amount"`
	Description        string    `json:"description"`
	ActiveStatus       bool      `json:"active_status"`
	AllowedCategoryIDs []string  `json:"allowed_category_ids"`
	AllowedProductIDs  []string  `json:"allowed_product_ids"`
	ExpiredAt          time.Time `json:"expired_at"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func GetVoucherTable() string {
	return "vouchers"
}
