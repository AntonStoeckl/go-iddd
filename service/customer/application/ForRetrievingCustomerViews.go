package application

import (
	"go-iddd/service/customer/application/readmodel/domain/customer"
)

type ForRetrievingCustomerViews func(customerID customer.ID) (customer.View, error)
