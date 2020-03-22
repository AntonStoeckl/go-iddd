package events

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	jsoniter "github.com/json-iterator/go"
)

type CustomerDeleted struct {
	customerID values.CustomerID
	meta       EventMeta
}

type CustomerDeletedForJSON struct {
	CustomerID string    `json:"customerID"`
	Meta       EventMeta `json:"meta"`
}

func CustomerWasDeleted(
	customerID values.CustomerID,
	streamVersion uint,
) CustomerDeleted {

	event := CustomerDeleted{
		customerID: customerID,
	}

	event.meta = BuildEventMeta(event, streamVersion)

	return event
}

func (event CustomerDeleted) CustomerID() values.CustomerID {
	return event.customerID
}

func (event CustomerDeleted) EventName() string {
	return event.meta.EventName
}

func (event CustomerDeleted) OccurredAt() string {
	return event.meta.OccurredAt
}

func (event CustomerDeleted) StreamVersion() uint {
	return event.meta.streamVersion
}

func (event CustomerDeleted) MarshalJSON() ([]byte, error) {
	data := CustomerDeletedForJSON{
		CustomerID: event.customerID.ID(),
		Meta:       event.meta,
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
		customerID: values.RebuildCustomerID(unmarshaledData.CustomerID),
		meta:       EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	}

	return event
}
