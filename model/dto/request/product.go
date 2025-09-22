package request

type CreateProductRequest struct {
	CategoryId   string  `json:"category_id" validate:"required"`
	CollectionId string  `json:"collection_id"`
	ProductName  string  `json:"collection_name"`
	Description  string  `json:"description"`
	Image        string  `json:"image"`
	Size         string  `json:"size"`
	Color        string  `json:"color"`
	Price        float64 `json:"price" validate:"required, min=0"`
	Currency     string  `json:"currency"`
	Amount       int64   `json:"amount"`
}

type UpdateProductRequest struct {
	ProductId    string  `json:"product_id" validate:"required"`
	CategoryId   string  `json:"category_id"`
	CollectionId string  `json:"collection_id"`
	ProductName  string  `json:"collection_name"`
	Description  string  `json:"description"`
	Image        string  `json:"image"`
	Size         string  `json:"size"`
	Color        string  `json:"color"`
	Price        float64 `json:"price"`
	Currency     string  `json:"currency"`
	Amount       *int64  `json:"amount"`
}

type GetProductsCustomerUI struct {
	Pagination   SearchPaginatioRequest `json:"pagination"`
	CategoryId   string                 `json:"category_id" form:"category_id"`
	CollectionId string                 `json:"collection_id" form:"collection_id"`
}

type ItemSelected struct {
	ProductId string `json:"product_id" validate:"required"`
	Quantity  int    `json:"quantity" validate:"requried,gt=0"`
}
