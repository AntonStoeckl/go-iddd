package hexagon

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer"
)

type ForRetrievingCustomerViews func(customerID string) (customer.View, error)
