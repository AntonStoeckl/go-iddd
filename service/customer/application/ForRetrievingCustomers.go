package application

import (
	"go-iddd/service/customer/application/readmodel"
)

type ForRetrievingCustomers func(query readmodel.CustomerByIDQuery) readmodel.CustomerView
