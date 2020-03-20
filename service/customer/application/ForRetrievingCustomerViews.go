package application

import (
	"go-iddd/service/customer/domain/customer"
	"go-iddd/service/customer/domain/customer/values"
)

type ForRetrievingCustomerViews func(customerID values.CustomerID) (customer.View, error)
