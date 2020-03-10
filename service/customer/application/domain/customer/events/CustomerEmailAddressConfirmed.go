package events

import (
	"go-iddd/service/customer/application/domain/customer/values"

	jsoniter "github.com/json-iterator/go"
)

type CustomerEmailAddressConfirmed struct {
	customerID   values.CustomerID
	emailAddress values.EmailAddress
	meta         EventMeta
}

func CustomerEmailAddressWasConfirmed(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	streamVersion uint,
) CustomerEmailAddressConfirmed {

	event := CustomerEmailAddressConfirmed{
		customerID:   customerID,
		emailAddress: emailAddress,
	}

	event.meta = BuildEventMeta(event, streamVersion)

	return event
}

func (event CustomerEmailAddressConfirmed) CustomerID() values.CustomerID {
	return event.customerID
}

func (event CustomerEmailAddressConfirmed) EmailAddress() values.EmailAddress {
	return event.emailAddress
}

func (event CustomerEmailAddressConfirmed) EventName() string {
	return event.meta.eventName
}

func (event CustomerEmailAddressConfirmed) OccurredAt() string {
	return event.meta.occurredAt
}

func (event CustomerEmailAddressConfirmed) StreamVersion() uint {
	return event.meta.streamVersion
}

func (event CustomerEmailAddressConfirmed) MarshalJSON() ([]byte, error) {
	data := &struct {
		CustomerID   string    `json:"customerID"`
		EmailAddress string    `json:"emailAddress"`
		Meta         EventMeta `json:"meta"`
	}{
		CustomerID:   event.customerID.ID(),
		EmailAddress: event.emailAddress.EmailAddress(),
		Meta:         event.meta,
	}

	return jsoniter.Marshal(data)
}

func UnmarshalCustomerEmailAddressConfirmedFromJSON(data []byte, streamVersion uint) CustomerEmailAddressConfirmed {
	event := CustomerEmailAddressConfirmed{
		customerID:   values.RebuildCustomerID(jsoniter.Get(data, "customerID").ToString()),
		emailAddress: values.RebuildEmailAddress(jsoniter.Get(data, "emailAddress").ToString()),
		meta:         UnmarshalEventMetaFromJSON(data, streamVersion),
	}

	return event
}
