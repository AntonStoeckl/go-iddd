package command

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type ForRetrievingCustomerEventStreams func(id values.CustomerID) (es.EventStream, error)
