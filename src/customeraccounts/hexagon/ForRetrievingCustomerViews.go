package hexagon

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer"
)

type ForRetrievingCustomerViews func(customerID string) (customer.View, error)
