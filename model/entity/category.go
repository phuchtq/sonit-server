package entity

import "time"

type Category struct {
	CategoryId   string    `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Description  string    `json:"description"`
	ActiveStatus bool      `json:"active_status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func GetCategoryTable() string {
	return "categories"
}
