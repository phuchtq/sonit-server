package entity

import "time"

type Product struct {
	ProductId    string    `json:"product_id"`
	CategoryId   string    `json:"category_id"`
	CollectionId string    `json:"collection_id"`
	ProductName  string    `json:"collection_name"`
	Description  string    `json:"description"`
	Image        string    `json:"image"`
	Size         string    `json:"size"`
	Color        string    `json:"color"`
	Price        float64   `json:"price"`
	Currency     string    `json:"currency"`
	ActiveStatus bool      `json:"active_status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func GetProductTable() string {
	return "products"
}
