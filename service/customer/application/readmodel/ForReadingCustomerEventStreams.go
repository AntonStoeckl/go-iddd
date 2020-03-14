package readmodel

import (
	"go-iddd/service/customer/application/readmodel/domain/customer"
	"go-iddd/service/lib/es"
)

type ForReadingCustomerEventStreams interface {
	EventStreamFor(id customer.ID) (es.DomainEvents, error)
}
