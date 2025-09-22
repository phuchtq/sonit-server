package entity

import "time"

type Shipping struct {
	OrderId        string     `json:"order_id"`
	DeliveryCode   string     `json:"delivery_code"`
	ShippingUnit   string     `json:"shipping_unit"`
	ShippingDetail string     `json:"shipping_detail"`
	DeliveredAt    *time.Time `json:"delivered_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func GetShippingTable() string {
	return "shippings"
}
