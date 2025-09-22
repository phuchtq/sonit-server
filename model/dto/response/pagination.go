package response

type PaginationDataResponse struct {
	Data       interface{} `json:"data"`
	PageNumber int         `json:"page_number"`
	TotalPages int         `json:"total_pages"`
}
