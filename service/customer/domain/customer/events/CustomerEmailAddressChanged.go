package events

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type CustomerEmailAddressChanged struct {
	customerID           values.CustomerID
	emailAddress         values.EmailAddress
	confirmationHash     values.ConfirmationHash
	previousEmailAddress values.EmailAddress
	meta                 es.EventMeta
}

func BuildCustomerEmailAddressChanged(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	confirmationHash values.ConfirmationHash,
	previousEmailAddress values.EmailAddress,
	streamVersion uint,
) CustomerEmailAddressChanged {

	event := CustomerEmailAddressChanged{
		customerID:           customerID,
		emailAddress:         emailAddress,
		confirmationHash:     confirmationHash,
		previousEmailAddress: previousEmailAddress,
	}

	event.meta = es.BuildEventMeta(event, streamVersion)

	return event
}

func RebuildCustomerEmailAddressChanged(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	confirmationHash values.ConfirmationHash,
	previousEmailAddress values.EmailAddress,
	meta es.EventMeta,
) CustomerEmailAddressChanged {

	event := CustomerEmailAddressChanged{
		customerID:           customerID,
		emailAddress:         emailAddress,
		confirmationHash:     confirmationHash,
		previousEmailAddress: previousEmailAddress,
		meta:                 meta,
	}

	return event
}

func (event CustomerEmailAddressChanged) CustomerID() values.CustomerID {
	return event.customerID
}

func (event CustomerEmailAddressChanged) EmailAddress() values.EmailAddress {
	return event.emailAddress
}

func (event CustomerEmailAddressChanged) ConfirmationHash() values.ConfirmationHash {
	return event.confirmationHash
}

func (event CustomerEmailAddressChanged) PreviousEmailAddress() values.EmailAddress {
	return event.previousEmailAddress
}

func (event CustomerEmailAddressChanged) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerEmailAddressChanged) IsFailureEvent() bool {
	return false
}

func (event CustomerEmailAddressChanged) FailureReason() error {
	return nil
}
