package events

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	jsoniter "github.com/json-iterator/go"
)

type CustomerDeleted struct {
	customerID   values.CustomerID
	emailAddress values.EmailAddress
	meta         EventMeta
}

type CustomerDeletedForJSON struct {
	CustomerID   string    `json:"customerID"`
	EmailAddress string    `json:"emailAddress"`
	Meta         EventMeta `json:"meta"`
}

func CustomerWasDeleted(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	streamVersion uint,
) CustomerDeleted {

	event := CustomerDeleted{
		customerID:   customerID,
		emailAddress: emailAddress,
	}

	event.meta = BuildEventMeta(event, streamVersion)

	return event
}

func (event CustomerDeleted) CustomerID() values.CustomerID {
	return event.customerID
}

func (event CustomerDeleted) EmailAddress() values.EmailAddress {
	return event.emailAddress
}

func (event CustomerDeleted) EventName() string {
	return event.meta.EventName
}

func (event CustomerDeleted) OccurredAt() string {
	return event.meta.OccurredAt
}

func (event CustomerDeleted) IndicatesAnError() (bool, string) {
	return false, ""
}

func (event CustomerDeleted) StreamVersion() uint {
	return event.meta.streamVersion
}

func (event CustomerDeleted) MarshalJSON() ([]byte, error) {
	data := CustomerDeletedForJSON{
		CustomerID:   event.customerID.String(),
		EmailAddress: event.emailAddress.String(),
		Meta:         event.meta,
	}

	return jsoniter.ConfigFastest.Marshal(data)
}

func UnmarshalCustomerDeletedFromJSON(
	data []byte,
	streamVersion uint,
) CustomerDeleted {

	unmarshaledData := &CustomerDeletedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData)

	event := CustomerDeleted{
		customerID:   values.RebuildCustomerID(unmarshaledData.CustomerID),
		emailAddress: values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		meta:         EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	}

	return event
}
