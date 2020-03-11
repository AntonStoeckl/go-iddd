package application

import (
	"go-iddd/service/customer/application/readmodel/domain/customer"
	"go-iddd/service/customer/application/readmodel/domain/customer/queries"
)

type ForRetrievingCustomerViews func(query queries.CustomerByID) (customer.View, error)
