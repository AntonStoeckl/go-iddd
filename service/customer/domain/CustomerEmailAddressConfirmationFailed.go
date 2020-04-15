package domain

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type CustomerEmailAddressConfirmationFailed struct {
	customerID       value.CustomerID
	emailAddress     value.EmailAddress
	confirmationHash value.ConfirmationHash
	reason           error
	meta             es.EventMeta
}

func BuildCustomerEmailAddressConfirmationFailed(
	customerID value.CustomerID,
	emailAddress value.EmailAddress,
	confirmationHash value.ConfirmationHash,
	reason error,
	streamVersion uint,
) CustomerEmailAddressConfirmationFailed {

	event := CustomerEmailAddressConfirmationFailed{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
		reason:           reason,
	}

	event.meta = es.BuildEventMeta(event, streamVersion)

	return event
}

func RebuildCustomerEmailAddressConfirmationFailed(
	customerID value.CustomerID,
	emailAddress value.EmailAddress,
	confirmationHash value.ConfirmationHash,
	reason error,
	meta es.EventMeta,
) CustomerEmailAddressConfirmationFailed {

	event := CustomerEmailAddressConfirmationFailed{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
		reason:           reason,
		meta:             meta,
	}

	return event
}

func (event CustomerEmailAddressConfirmationFailed) CustomerID() value.CustomerID {
	return event.customerID
}

func (event CustomerEmailAddressConfirmationFailed) EmailAddress() value.EmailAddress {
	return event.emailAddress
}

func (event CustomerEmailAddressConfirmationFailed) ConfirmationHash() value.ConfirmationHash {
	return event.confirmationHash
}

func (event CustomerEmailAddressConfirmationFailed) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerEmailAddressConfirmationFailed) IsFailureEvent() bool {
	return true
}

func (event CustomerEmailAddressConfirmationFailed) FailureReason() error {
	return event.reason
}
