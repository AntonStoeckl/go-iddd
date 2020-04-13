package command

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type ForStoringCustomerEvents interface {
	RetrieveCustomerEventStream(id values.CustomerID) (es.DomainEvents, error)
	RegisterCustomer(recordedEvents es.DomainEvents, id values.CustomerID) error
	AppendToCustomerEventStream(recordedEvents es.DomainEvents, id values.CustomerID) error
	PurgeCustomerEventStream(id values.CustomerID) error
}
