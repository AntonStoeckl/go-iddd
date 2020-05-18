package domain

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/shared/es"
)

type CustomerEmailAddressChanged struct {
	customerID           value.CustomerID
	emailAddress         value.EmailAddress
	confirmationHash     value.ConfirmationHash
	previousEmailAddress value.EmailAddress
	meta                 es.EventMeta
}

func BuildCustomerEmailAddressChanged(
	customerID value.CustomerID,
	emailAddress value.EmailAddress,
	confirmationHash value.ConfirmationHash,
	previousEmailAddress value.EmailAddress,
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
	customerID string,
	emailAddress string,
	confirmationHash string,
	previousEmailAddress string,
	meta es.EventMeta,
) CustomerEmailAddressChanged {

	event := CustomerEmailAddressChanged{
		customerID:           value.RebuildCustomerID(customerID),
		emailAddress:         value.RebuildEmailAddress(emailAddress),
		confirmationHash:     value.RebuildConfirmationHash(confirmationHash),
		previousEmailAddress: value.RebuildEmailAddress(previousEmailAddress),
		meta:                 meta,
	}

	return event
}

func (event CustomerEmailAddressChanged) CustomerID() value.CustomerID {
	return event.customerID
}

func (event CustomerEmailAddressChanged) EmailAddress() value.EmailAddress {
	return event.emailAddress
}

func (event CustomerEmailAddressChanged) ConfirmationHash() value.ConfirmationHash {
	return event.confirmationHash
}

func (event CustomerEmailAddressChanged) PreviousEmailAddress() value.EmailAddress {
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
