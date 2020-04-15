package domain

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type CustomerEmailAddressConfirmed struct {
	customerID   value.CustomerID
	emailAddress value.EmailAddress
	meta         es.EventMeta
}

func BuildCustomerEmailAddressConfirmed(
	customerID value.CustomerID,
	emailAddress value.EmailAddress,
	streamVersion uint,
) CustomerEmailAddressConfirmed {

	event := CustomerEmailAddressConfirmed{
		customerID:   customerID,
		emailAddress: emailAddress,
	}

	event.meta = es.BuildEventMeta(event, streamVersion)

	return event
}

func RebuildCustomerEmailAddressConfirmed(
	customerID value.CustomerID,
	emailAddress value.EmailAddress,
	meta es.EventMeta,
) CustomerEmailAddressConfirmed {

	event := CustomerEmailAddressConfirmed{
		customerID:   customerID,
		emailAddress: emailAddress,
		meta:         meta,
	}

	return event
}

func (event CustomerEmailAddressConfirmed) CustomerID() value.CustomerID {
	return event.customerID
}

func (event CustomerEmailAddressConfirmed) EmailAddress() value.EmailAddress {
	return event.emailAddress
}

func (event CustomerEmailAddressConfirmed) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerEmailAddressConfirmed) IsFailureEvent() bool {
	return false
}

func (event CustomerEmailAddressConfirmed) FailureReason() error {
	return nil
}
