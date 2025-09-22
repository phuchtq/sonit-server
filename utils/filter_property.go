package utils

import (
	filter_property "sonit_server/constant/filter_property"
)

func AssignFilterProperty(filterProp string) string {
	var res string

	switch filterProp {
	case filter_property.DATE_FILTER:
		res = "created_at"
	case filter_property.ACTION_DATE_FILTER:
		res = "date"
	case filter_property.PRICE_FILTER:
		res = "price"
	case filter_property.RATE_FILTER:
		res = "rate"
	case filter_property.AMOUNT_FILTER:
		res = "amount"
	default:
		res = "created_at"
	}

	return res
}
