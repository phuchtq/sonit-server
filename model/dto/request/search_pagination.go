package request

type SearchPaginatioRequest struct {
	PageNumber int    `json:"page_number" form:"page_number"`
	Keyword    string `json:"keyword" form:"keyword"`
	FilterProp string `json:"filter_prop" form:"filter_prop"` // Date, price, ...
	Order      string `json:"order" form:"order"`             //ASC or DESC
}
