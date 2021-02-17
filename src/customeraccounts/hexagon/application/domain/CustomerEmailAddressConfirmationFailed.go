package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
)

type CustomerEmailAddressConfirmationFailed struct {
	customerID       value.CustomerID
	confirmationHash value.ConfirmationHash
	reason           error
	meta             es.EventMeta
}

func BuildCustomerEmailAddressConfirmationFailed(
	customerID value.CustomerID,
	confirmationHash value.ConfirmationHash,
	reason error,
	causationID es.MessageID,
	streamVersion uint,
) CustomerEmailAddressConfirmationFailed {

	event := CustomerEmailAddressConfirmationFailed{
		customerID:       customerID,
		confirmationHash: confirmationHash,
		reason:           reason,
	}

	event.meta = es.BuildEventMeta(event, causationID, streamVersion)

	return event
}

func RebuildCustomerEmailAddressConfirmationFailed(
	customerID string,
	confirmationHash string,
	reason string,
	meta es.EventMeta,
) CustomerEmailAddressConfirmationFailed {

	event := CustomerEmailAddressConfirmationFailed{
		customerID:       value.RebuildCustomerID(customerID),
		confirmationHash: value.RebuildConfirmationHash(confirmationHash),
		reason:           errors.Mark(errors.New(reason), shared.ErrDomainConstraintsViolation),
		meta:             meta,
	}

	return event
}

func (event CustomerEmailAddressConfirmationFailed) CustomerID() value.CustomerID {
	return event.customerID
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
