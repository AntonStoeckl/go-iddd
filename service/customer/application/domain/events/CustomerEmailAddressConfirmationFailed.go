package events

import (
	"go-iddd/service/customer/application/domain/values"

	jsoniter "github.com/json-iterator/go"
)

type CustomerEmailAddressConfirmationFailed struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	meta             EventMeta
}

func CustomerEmailAddressConfirmationHasFailed(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	confirmationHash values.ConfirmationHash,
	streamVersion uint,
) CustomerEmailAddressConfirmationFailed {

	event := CustomerEmailAddressConfirmationFailed{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
	}

	event.meta = BuildEventMeta(event, streamVersion)

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

func (event CustomerEmailAddressConfirmationFailed) EventName() string {
	return event.meta.eventName
}

func (event CustomerEmailAddressConfirmationFailed) OccurredAt() string {
	return event.meta.occurredAt
}

func (event CustomerEmailAddressConfirmationFailed) StreamVersion() uint {
	return event.meta.streamVersion
}

func (event CustomerEmailAddressConfirmationFailed) MarshalJSON() ([]byte, error) {
	data := struct {
		CustomerID       string    `json:"customerID"`
		EmailAddress     string    `json:"emailAddress"`
		ConfirmationHash string    `json:"confirmationHash"`
		Meta             EventMeta `json:"meta"`
	}{
		CustomerID:       event.customerID.ID(),
		EmailAddress:     event.emailAddress.EmailAddress(),
		ConfirmationHash: event.confirmationHash.Hash(),
		Meta:             event.meta,
	}

	return jsoniter.Marshal(data)
}

func UnmarshalCustomerEmailAddressConfirmationFailedFromJSON(data []byte, streamVersion uint) CustomerEmailAddressConfirmationFailed {
	event := CustomerEmailAddressConfirmationFailed{
		customerID:       values.RebuildCustomerID(jsoniter.Get(data, "customerID").ToString()),
		emailAddress:     values.RebuildEmailAddress(jsoniter.Get(data, "emailAddress").ToString()),
		confirmationHash: values.RebuildConfirmationHash(jsoniter.Get(data, "confirmationHash").ToString()),
		meta:             UnmarshalEventMetaFromJSON(data, streamVersion),
	}

	return event
}
