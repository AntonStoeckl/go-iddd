package customeraccounts

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer"
)

type ForRetrievingCustomerViews func(customerID string) (customer.View, error)
