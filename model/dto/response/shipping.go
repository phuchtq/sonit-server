package response

import "time"

type ShippingDetail struct {
	RecipientName string `json:"recipient_name"`
	Address       string `json:"address"`
	City          string `json:"city"`
	Country       string `json:"country"`
	PhoneNumber   string `json:"phone_number"`
}

type ShippingResponse struct {
	OrderId        string         `json:"order_id"`
	DeliveryCode   string         `json:"delivery_code"`
	ShippingUnit   string         `json:"shipping_unit"`
	ShippingDetail ShippingDetail `json:"shipping_detail"`
	DeliveredAt    *time.Time     `json:"delivered_at"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}
