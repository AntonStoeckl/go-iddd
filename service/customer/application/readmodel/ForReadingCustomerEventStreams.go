package readmodel

import (
	"go-iddd/service/lib/es"
)

type ForReadingCustomerEventStreams interface {
	EventStreamFor(id es.AggregateID) (es.DomainEvents, error)
}
