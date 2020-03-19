package application

import (
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib/es"
)

type ForReadingCustomerEventStreams interface {
	EventStreamFor(id values.CustomerID) (es.DomainEvents, error)
}
