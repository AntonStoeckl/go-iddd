package application

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type ForRetrievingCustomerEventStreams func(id value.CustomerID) (es.EventStream, error)
