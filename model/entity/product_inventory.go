package entity

type ProductInventory struct {
	ProductId       string `json:"product_id"`
	CurrentQuantity int64  `json:"current_quantity"`
}

func GetProductInventoryTable() string {
	return "product_inventories"
}
