package writemodel

import (
	"go-iddd/service/customer/application/writemodel/domain/customer/values"
	"go-iddd/service/lib/es"
)

type ForStoringCustomerEvents interface {
	EventStreamFor(id es.AggregateID) (es.DomainEvents, error)
	CreateStreamFrom(recordedEvents es.DomainEvents, id values.CustomerID) error
	Add(recordedEvents es.DomainEvents, id values.CustomerID) error
	Delete(id values.CustomerID) error
}
