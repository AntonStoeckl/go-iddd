package events

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type CustomerEmailAddressConfirmationFailed struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	reason           string
	meta             es.EventMeta
}

func BuildCustomerEmailAddressConfirmationFailed(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	confirmationHash values.ConfirmationHash,
	reason string,
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

func (event CustomerEmailAddressConfirmationFailed) CustomerID() values.CustomerID {
	return event.customerID
}

func (event CustomerEmailAddressConfirmationFailed) EmailAddress() values.EmailAddress {
	return event.emailAddress
}

func (event CustomerEmailAddressConfirmationFailed) ConfirmationHash() values.ConfirmationHash {
	return event.confirmationHash
}

func (event CustomerEmailAddressConfirmationFailed) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerEmailAddressConfirmationFailed) EventName() string {
	return event.meta.EventName
}

func (event CustomerEmailAddressConfirmationFailed) OccurredAt() string {
	return event.meta.OccurredAt
}

func (event CustomerEmailAddressConfirmationFailed) StreamVersion() uint {
	return event.meta.StreamVersion
}

func (event CustomerEmailAddressConfirmationFailed) IndicatesAnError() (bool, string) {
	return true, event.reason
}
