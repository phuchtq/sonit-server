package entity

import "time"

type Collection struct {
	CollectionId   string    `json:"collection_id"`
	CollectionName string    `json:"collection_name"`
	Description    string    `json:"description"`
	ActiveStatus   bool      `json:"active_status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func GetCollectionTable() string {
	return "collections"
}
