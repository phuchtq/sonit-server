package entity

import "time"

type ProductInventoryTransaction struct {
	TransactionId string    `json:"transaction_id"`
	ProductId     string    `json:"product_id"`
	Amount        int64     `json:"amount"`
	Action        string    `json:"action"` // Import, Export, Sale
	Note          string    `json:"note"`
	Date          time.Time `json:"date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func GetProductInventoryTransactionTable() string {
	return "product_inventory_transactions"
}
