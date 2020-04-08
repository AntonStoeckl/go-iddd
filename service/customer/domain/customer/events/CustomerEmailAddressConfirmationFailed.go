package events

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	jsoniter "github.com/json-iterator/go"
)

type CustomerEmailAddressConfirmationFailed struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	reason           string
	meta             es.EventMeta
}

type CustomerEmailAddressConfirmationFailedForJSON struct {
	CustomerID       string       `json:"customerID"`
	EmailAddress     string       `json:"emailAddress"`
	ConfirmationHash string       `json:"confirmationHash"`
	Reason           string       `json:"reason"`
	Meta             es.EventMeta `json:"meta"`
}

func CustomerEmailAddressConfirmationHasFailed(
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

func (event CustomerEmailAddressConfirmationFailed) MarshalJSON() ([]byte, error) {
	data := CustomerEmailAddressConfirmationFailedForJSON{
		CustomerID:       event.customerID.String(),
		EmailAddress:     event.emailAddress.String(),
		ConfirmationHash: event.confirmationHash.String(),
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
		meta:             es.EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	}

	return event
}
