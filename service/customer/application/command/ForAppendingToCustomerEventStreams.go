package command

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type ForAppendingToCustomerEventStreams func(recordedEvents es.DomainEvents, id values.CustomerID) error
