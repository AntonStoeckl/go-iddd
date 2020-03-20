package command

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type ForStoringCustomerEvents interface {
	EventStreamFor(id values.CustomerID) (es.DomainEvents, error)
	CreateStreamFrom(recordedEvents es.DomainEvents, id values.CustomerID) error
	Add(recordedEvents es.DomainEvents, id values.CustomerID) error
	Delete(id values.CustomerID) error
}
