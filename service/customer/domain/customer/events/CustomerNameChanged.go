package events

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type CustomerNameChanged struct {
	customerID values.CustomerID
	personName values.PersonName
	meta       es.EventMeta
}

func BuildCustomerNameChanged(
	customerID values.CustomerID,
	personName values.PersonName,
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
	customerID values.CustomerID,
	personName values.PersonName,
	meta es.EventMeta,
) CustomerNameChanged {

	event := CustomerNameChanged{
		customerID: customerID,
		personName: personName,
		meta:       meta,
	}

	return event
}

func (event CustomerNameChanged) CustomerID() values.CustomerID {
	return event.customerID
}

func (event CustomerNameChanged) PersonName() values.PersonName {
	return event.personName
}

func (event CustomerNameChanged) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerNameChanged) EventName() string {
	return event.meta.EventName
}

func (event CustomerNameChanged) OccurredAt() string {
	return event.meta.OccurredAt
}

func (event CustomerNameChanged) StreamVersion() uint {
	return event.meta.StreamVersion
}

func (event CustomerNameChanged) IndicatesAnError() (bool, string) {
	return false, ""
}
