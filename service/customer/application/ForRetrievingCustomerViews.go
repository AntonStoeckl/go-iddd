package application

import (
	"go-iddd/service/customer/application/domain/customer"
	"go-iddd/service/customer/application/domain/values"
)

type ForRetrievingCustomerViews func(customerID values.CustomerID) (customer.View, error)
