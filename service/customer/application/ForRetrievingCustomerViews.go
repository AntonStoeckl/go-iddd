package application

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer"
)

type ForRetrievingCustomerViews func(customerID string) (customer.View, error)
