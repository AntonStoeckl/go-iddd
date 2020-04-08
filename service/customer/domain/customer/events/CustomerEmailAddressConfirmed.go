package events

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	jsoniter "github.com/json-iterator/go"
)

type CustomerEmailAddressConfirmed struct {
	customerID   values.CustomerID
	emailAddress values.EmailAddress
	meta         es.EventMeta
}

type CustomerEmailAddressConfirmedForJSON struct {
	CustomerID   string       `json:"customerID"`
	EmailAddress string       `json:"emailAddress"`
	Meta         es.EventMeta `json:"meta"`
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

	event.meta = es.BuildEventMeta(event, streamVersion)

	return event
}

func (event CustomerEmailAddressConfirmed) CustomerID() values.CustomerID {
	return event.customerID
}

func (event CustomerEmailAddressConfirmed) EmailAddress() values.EmailAddress {
	return event.emailAddress
}

func (event CustomerEmailAddressConfirmed) EventName() string {
	return event.meta.EventName
}

func (event CustomerEmailAddressConfirmed) OccurredAt() string {
	return event.meta.OccurredAt
}

func (event CustomerEmailAddressConfirmed) StreamVersion() uint {
	return event.meta.StreamVersion
}

func (event CustomerEmailAddressConfirmed) IndicatesAnError() (bool, string) {
	return false, ""
}

func (event CustomerEmailAddressConfirmed) MarshalJSON() ([]byte, error) {
	data := CustomerEmailAddressConfirmedForJSON{
		CustomerID:   event.customerID.String(),
		EmailAddress: event.emailAddress.String(),
		Meta:         event.meta,
	}

	return jsoniter.ConfigFastest.Marshal(data)
}

func UnmarshalCustomerEmailAddressConfirmedFromJSON(
	data []byte,
	streamVersion uint,
) CustomerEmailAddressConfirmed {

	unmarshaledData := &CustomerEmailAddressConfirmedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData)

	event := CustomerEmailAddressConfirmed{
		customerID:   values.RebuildCustomerID(unmarshaledData.CustomerID),
		emailAddress: values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		meta:         es.EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	}

	return event
}
