package events

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type CustomerDeleted struct {
	customerID   values.CustomerID
	emailAddress values.EmailAddress
	meta         es.EventMeta
}

func BuildCustomerDeleted(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	streamVersion uint,
) CustomerDeleted {

	event := CustomerDeleted{
		customerID:   customerID,
		emailAddress: emailAddress,
	}

	event.meta = es.BuildEventMeta(event, streamVersion)

	return event
}

func RebuildCustomerDeleted(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	meta es.EventMeta,
) CustomerDeleted {

	event := CustomerDeleted{
		customerID:   customerID,
		emailAddress: emailAddress,
		meta:         meta,
	}

	return event
}

func (event CustomerDeleted) CustomerID() values.CustomerID {
	return event.customerID
}

func (event CustomerDeleted) EmailAddress() values.EmailAddress {
	return event.emailAddress
}

func (event CustomerDeleted) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerDeleted) IndicatesAnError() (bool, string) {
	return false, ""
}
