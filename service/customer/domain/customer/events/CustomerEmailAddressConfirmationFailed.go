package events

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	jsoniter "github.com/json-iterator/go"
)

type CustomerEmailAddressConfirmationFailed struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	meta             EventMeta
}

type CustomerEmailAddressConfirmationFailedForJSON struct {
	CustomerID       string    `json:"customerID"`
	EmailAddress     string    `json:"emailAddress"`
	ConfirmationHash string    `json:"confirmationHash"`
	Meta             EventMeta `json:"meta"`
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
	return event.meta.EventName
}

func (event CustomerEmailAddressConfirmationFailed) OccurredAt() string {
	return event.meta.OccurredAt
}

func (event CustomerEmailAddressConfirmationFailed) StreamVersion() uint {
	return event.meta.streamVersion
}

func (event CustomerEmailAddressConfirmationFailed) MarshalJSON() ([]byte, error) {
	data := CustomerEmailAddressConfirmationFailedForJSON{
		CustomerID:       event.customerID.ID(),
		EmailAddress:     event.emailAddress.EmailAddress(),
		ConfirmationHash: event.confirmationHash.Hash(),
		Meta:             event.meta,
	}

	return jsoniter.ConfigFastest.Marshal(data)
}

func UnmarshalCustomerEmailAddressConfirmationFailedFromJSON(
	data []byte,
	streamVersion uint,
) CustomerEmailAddressConfirmationFailed {

	unmarshaledData := &CustomerEmailAddressConfirmationFailedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData)

	event := CustomerEmailAddressConfirmationFailed{
		customerID:       values.RebuildCustomerID(unmarshaledData.CustomerID),
		emailAddress:     values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		confirmationHash: values.RebuildConfirmationHash(unmarshaledData.ConfirmationHash),
		meta:             EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	}

	return event
}
