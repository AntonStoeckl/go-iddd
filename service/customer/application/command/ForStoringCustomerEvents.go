package command

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type ForStoringCustomerEvents interface {
	RetrieveCustomerEventStream(id values.CustomerID) (es.EventStream, error)
	RegisterCustomer(recordedEvents es.RecordedEvents, id values.CustomerID) error
	AppendToCustomerEventStream(recordedEvents es.RecordedEvents, id values.CustomerID) error
	PurgeCustomerEventStream(id values.CustomerID) error
}
