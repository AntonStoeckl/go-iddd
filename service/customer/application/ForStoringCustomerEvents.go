package application

import (
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib/es"
)

type ForStoringCustomerEvents interface {
	EventStreamFor(id values.CustomerID) (es.DomainEvents, error)
	CreateStreamFrom(recordedEvents es.DomainEvents, id values.CustomerID) error
	Add(recordedEvents es.DomainEvents, id values.CustomerID) error
	Delete(id values.CustomerID) error
}
