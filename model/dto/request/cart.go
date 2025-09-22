package request

type AddItemToCartRequest struct {
	Request  RemoveItemFromCartRequest `json:"request" validate:"required"`
	Quantity int                       `json:"quantity" validate:"required, min=1"`
}

type RemoveItemFromCartRequest struct {
	UserId    string `json:"user_id" validate:"required"`
	ProductId string `json:"product_id" validate:"required"`
}

type EditItemInCartRequest struct {
	Request  RemoveItemFromCartRequest `json:"request"`
	Quantity int                       `json:"quantity" validate:"required, min=1"`
}
