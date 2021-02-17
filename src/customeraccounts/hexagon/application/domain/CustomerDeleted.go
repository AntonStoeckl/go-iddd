package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type CustomerDeleted struct {
	customerID value.CustomerID
	meta       es.EventMeta
}

func BuildCustomerDeleted(
	customerID value.CustomerID,
	causationID es.MessageID,
	streamVersion uint,
) CustomerDeleted {

	event := CustomerDeleted{
		customerID: customerID,
	}

	event.meta = es.BuildEventMeta(event, causationID, streamVersion)

	return event
}

func RebuildCustomerDeleted(
	customerID string,
	meta es.EventMeta,
) CustomerDeleted {

	event := CustomerDeleted{
		customerID: value.RebuildCustomerID(customerID),
		meta:       meta,
	}

	return event
}

func (event CustomerDeleted) CustomerID() value.CustomerID {
	return event.customerID
}

func (event CustomerDeleted) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerDeleted) IsFailureEvent() bool {
	return false
}

func (event CustomerDeleted) FailureReason() error {
	return nil
}
