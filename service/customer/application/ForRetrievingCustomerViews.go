package application

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
)

type ForRetrievingCustomerViews func(customerID values.CustomerID) (customer.View, error)
