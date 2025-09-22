package request

import "time"

type GetProductInventoryTrasactionsRequest struct {
	Pagination SearchPaginatioRequest `json:"pagination" validate:"required"`
	Action     string                 `json:"action" form:"action"`
}

type GetProductInventoryTrasactionsByProductRequest struct {
	InventoryTransaction GetProductInventoryTrasactionsRequest `json:"inventory_transaction" validate:"required"`
	ProductId            string                                `json:"product_id" form:"product_id" validate:"required"`
}

type CreateProductInventoryTransactionRequest struct {
	ProductId string    `json:"product_id" validate:"required"`
	Amount    int64     `json:"amount" validate:"required, min=1"`
	Action    string    `json:"action" validate:"required"` // Import, Export, Sale
	Note      string    `json:"note"`
	Date      time.Time `json:"date" validate:"required"`
}

type UpdateProductInventoryTransactionRequest struct {
	TransactionId string                                   `json:"transaction_id" validate:"required"`
	UpdateData    CreateProductInventoryTransactionRequest `json:"update_data"`
}
