package request

type CreateCategoryRequest struct {
	CategoryName string `json:"category_name" validate:"required"`
	Description  string `json:"description"`
}

type UpdateCategoryRequest struct {
	CategoryId   string `json:"category_id" validate:"required"`
	CategoryName string `json:"category_name"`
	Description  string `json:"description"`
}
