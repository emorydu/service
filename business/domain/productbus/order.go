package productbus

import "github.com/emorydu/service/business/sdk/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByProductID, order.DESC)

// Set of fields that the results can be ordered by.
const (
	OrderByProductID = "a"
	OrderByUserID    = "b"
	OrderByName      = "c"
	OrderByCost      = "d"
	OrderByQuantity  = "e"
)
