package readmodel

import (
	"go-iddd/service/customer/application/readmodel/domain/customer/values"
	"go-iddd/service/lib/es"
)

type ForReadingCustomerEventStreams interface {
	EventStreamFor(id values.CustomerID) (es.DomainEvents, error)
}
