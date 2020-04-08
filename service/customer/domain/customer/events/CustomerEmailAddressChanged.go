package events

import (
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	jsoniter "github.com/json-iterator/go"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
)

type CustomerEmailAddressChanged struct {
	customerID           values.CustomerID
	emailAddress         values.EmailAddress
	confirmationHash     values.ConfirmationHash
	previousEmailAddress values.EmailAddress
	meta                 es.EventMeta
}

type CustomerEmailAddressChangedForJSON struct {
	CustomerID           string       `json:"customerID"`
	EmailAddress         string       `json:"emailAddress"`
	ConfirmationHash     string       `json:"confirmationHash"`
	PreviousEmailAddress string       `json:"previousEmailAddress"`
	Meta                 es.EventMeta `json:"meta"`
}

func CustomerEmailAddressWasChanged(
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

func (event CustomerEmailAddressChanged) EventName() string {
	return event.meta.EventName
}

func (event CustomerEmailAddressChanged) OccurredAt() string {
	return event.meta.OccurredAt
}

func (event CustomerEmailAddressChanged) StreamVersion() uint {
	return event.meta.StreamVersion
}

func (event CustomerEmailAddressChanged) IndicatesAnError() (bool, string) {
	return false, ""
}

func (event CustomerEmailAddressChanged) MarshalJSON() ([]byte, error) {
	data := CustomerEmailAddressChangedForJSON{
		CustomerID:           event.customerID.String(),
		EmailAddress:         event.emailAddress.String(),
		ConfirmationHash:     event.confirmationHash.String(),
		PreviousEmailAddress: event.previousEmailAddress.String(),
		Meta:                 event.meta,
	}

	return jsoniter.ConfigFastest.Marshal(data)
}

func UnmarshalCustomerEmailAddressChangedFromJSON(
	data []byte,
	streamVersion uint,
) CustomerEmailAddressChanged {

	unmarshaledData := &CustomerEmailAddressChangedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData)

	event := CustomerEmailAddressChanged{
		customerID:           values.RebuildCustomerID(unmarshaledData.CustomerID),
		emailAddress:         values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		confirmationHash:     values.RebuildConfirmationHash(unmarshaledData.ConfirmationHash),
		previousEmailAddress: values.RebuildEmailAddress(unmarshaledData.PreviousEmailAddress),
		meta:                 es.EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	}

	return event
}
