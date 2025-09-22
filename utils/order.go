package utils

import "sonit_server/constant/order"

func AssignOrder(ord string) string {
	var res string = ord

	switch res {
	case order.ASCENDING_ORDER:
	case order.DESCENDING_ORDER:
	default:
		res = order.ASCENDING_ORDER
	}

	return res
}
