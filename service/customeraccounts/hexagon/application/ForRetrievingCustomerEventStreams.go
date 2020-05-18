package application

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/shared/es"
)

type ForRetrievingCustomerEventStreams func(id value.CustomerID) (es.EventStream, error)
