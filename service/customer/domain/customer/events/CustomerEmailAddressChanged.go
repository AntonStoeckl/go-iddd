package events

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	jsoniter "github.com/json-iterator/go"
)

type CustomerEmailAddressChanged struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	meta             EventMeta
}

func CustomerEmailAddressWasChanged(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	confirmationHash values.ConfirmationHash,
	streamVersion uint,
) CustomerEmailAddressChanged {

	event := CustomerEmailAddressChanged{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
	}

	event.meta = BuildEventMeta(event, streamVersion)

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

func (event CustomerEmailAddressChanged) EventName() string {
	return event.meta.eventName
}

func (event CustomerEmailAddressChanged) OccurredAt() string {
	return event.meta.occurredAt
}

func (event CustomerEmailAddressChanged) StreamVersion() uint {
	return event.meta.streamVersion
}

func (event CustomerEmailAddressChanged) MarshalJSON() ([]byte, error) {
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

	return jsoniter.ConfigFastest.Marshal(data)
}

func UnmarshalCustomerEmailAddressChangedFromJSON(
	data []byte,
	streamVersion uint,
) CustomerEmailAddressChanged {

	anyData := jsoniter.ConfigFastest.Get(data)

	event := CustomerEmailAddressChanged{
		customerID:       values.RebuildCustomerID(anyData.Get("customerID").ToString()),
		emailAddress:     values.RebuildEmailAddress(anyData.Get("emailAddress").ToString()),
		confirmationHash: values.RebuildConfirmationHash(anyData.Get("confirmationHash").ToString()),
		meta:             UnmarshalEventMetaFromJSON(data, streamVersion),
	}

	return event
}
