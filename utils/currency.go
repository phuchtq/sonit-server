package utils

import (
	"sonit_server/constant/currency"
)

func AssignCurrency(curr string) string {
	var res string = curr

	switch res {
	case currency.US_DOLLAR:
	case currency.VIETNAM_DONG:
	default:
		res = currency.VIETNAM_DONG
	}

	return res
}
