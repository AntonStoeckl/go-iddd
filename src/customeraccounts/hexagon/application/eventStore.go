package application

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type EventStoreInterface interface {
	RetrieveEventStream(id value.CustomerID) (es.EventStream, error)
	StartEventStream(customerRegistered domain.CustomerRegistered) error
	AppendToEventStream(recordedEvents es.RecordedEvents, id value.CustomerID) error
	PurgeEventStream(id value.CustomerID) error
}
