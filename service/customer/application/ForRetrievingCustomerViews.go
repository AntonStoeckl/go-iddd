package application

import (
	"go-iddd/service/customer/application/readmodel/domain/customer"
	"go-iddd/service/customer/application/readmodel/domain/customer/values"
)

type ForRetrievingCustomerViews func(customerID values.CustomerID) (customer.View, error)
