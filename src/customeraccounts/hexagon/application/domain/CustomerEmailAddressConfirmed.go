package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type CustomerEmailAddressConfirmed struct {
	customerID   value.CustomerID
	emailAddress value.EmailAddress
	meta         es.EventMeta
}

func BuildCustomerEmailAddressConfirmed(
	customerID value.CustomerID,
	emailAddress value.EmailAddress,
	causationID es.MessageID,
	streamVersion uint,
) CustomerEmailAddressConfirmed {

	event := CustomerEmailAddressConfirmed{
		customerID:   customerID,
		emailAddress: emailAddress,
	}

	event.meta = es.BuildEventMeta(event, causationID, streamVersion)

	return event
}

func RebuildCustomerEmailAddressConfirmed(
	customerID string,
	emailAddress string,
	meta es.EventMeta,
) CustomerEmailAddressConfirmed {

	event := CustomerEmailAddressConfirmed{
		customerID:   value.RebuildCustomerID(customerID),
		emailAddress: value.RebuildEmailAddress(emailAddress),
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
