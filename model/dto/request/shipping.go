package request

type GetShippingsRequest struct {
	Request SearchPaginatioRequest `json:"request"`
	UserId  string                 `json:"user_id"`
}
