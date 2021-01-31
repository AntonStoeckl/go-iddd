package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
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
	customerID string,
	emailAddress string,
	confirmationHash string,
	reason string,
	meta es.EventMeta,
) CustomerEmailAddressConfirmationFailed {

	event := CustomerEmailAddressConfirmationFailed{
		customerID:       value.RebuildCustomerID(customerID),
		emailAddress:     value.RebuildEmailAddress(emailAddress),
		confirmationHash: value.RebuildConfirmationHash(confirmationHash),
		reason:           errors.Mark(errors.New(reason), shared.ErrDomainConstraintsViolation),
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
