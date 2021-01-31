package application

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type ForRetrievingCustomerEventStreams func(id value.CustomerID) (es.EventStream, error)
