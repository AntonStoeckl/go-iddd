package domain

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/shared/es"
)

type CustomerNameChanged struct {
	customerID value.CustomerID
	personName value.PersonName
	meta       es.EventMeta
}

func BuildCustomerNameChanged(
	customerID value.CustomerID,
	personName value.PersonName,
	streamVersion uint,
) CustomerNameChanged {

	event := CustomerNameChanged{
		customerID: customerID,
		personName: personName,
	}

	event.meta = es.BuildEventMeta(event, streamVersion)

	return event
}

func RebuildCustomerNameChanged(
	customerID string,
	givenName string,
	familyName string,
	meta es.EventMeta,
) CustomerNameChanged {

	event := CustomerNameChanged{
		customerID: value.RebuildCustomerID(customerID),
		personName: value.RebuildPersonName(givenName, familyName),
		meta:       meta,
	}

	return event
}

func (event CustomerNameChanged) CustomerID() value.CustomerID {
	return event.customerID
}

func (event CustomerNameChanged) PersonName() value.PersonName {
	return event.personName
}

func (event CustomerNameChanged) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerNameChanged) IsFailureEvent() bool {
	return false
}

func (event CustomerNameChanged) FailureReason() error {
	return nil
}
