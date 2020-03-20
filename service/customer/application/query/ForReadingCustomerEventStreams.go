package query

import (
	"go-iddd/service/customer/domain/customer/values"
	"go-iddd/service/lib/es"
)

type ForReadingCustomerEventStreams interface {
	EventStreamFor(id values.CustomerID) (es.DomainEvents, error)
}
