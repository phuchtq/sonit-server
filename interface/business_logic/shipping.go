package businesslogic

import (
	"context"
	"sonit_server/model/dto/request"
	"sonit_server/model/dto/response"
)

type IShippingService interface {
	GetShippings(req request.GetShippingsRequest, ctx context.Context) (response.PaginationDataResponse, error)
}
