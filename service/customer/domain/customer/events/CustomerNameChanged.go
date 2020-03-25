package events

import (
	jsoniter "github.com/json-iterator/go"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
)

type CustomerNameChanged struct {
	customerID values.CustomerID
	personName values.PersonName
	meta       EventMeta
}

type CustomerNameChangedForJSON struct {
	CustomerID string    `json:"customerID"`
	GivenName  string    `json:"givenName"`
	FamilyName string    `json:"familyName"`
	Meta       EventMeta `json:"meta"`
}

func CustomerNameWasChanged(
	customerID values.CustomerID,
	personName values.PersonName,
	streamVersion uint,
) CustomerNameChanged {

	event := CustomerNameChanged{
		customerID: customerID,
		personName: personName,
	}

	event.meta = BuildEventMeta(event, streamVersion)

	return event
}

func (event CustomerNameChanged) CustomerID() values.CustomerID {
	return event.customerID
}

func (event CustomerNameChanged) PersonName() values.PersonName {
	return event.personName
}

func (event CustomerNameChanged) EventName() string {
	return event.meta.EventName
}

func (event CustomerNameChanged) OccurredAt() string {
	return event.meta.OccurredAt
}

func (event CustomerNameChanged) IndicatesAnError() (bool, string) {
	return false, ""
}

func (event CustomerNameChanged) StreamVersion() uint {
	return event.meta.streamVersion
}

func (event CustomerNameChanged) MarshalJSON() ([]byte, error) {
	data := CustomerNameChangedForJSON{
		CustomerID: event.customerID.ID(),
		GivenName:  event.personName.GivenName(),
		FamilyName: event.personName.FamilyName(),
		Meta:       event.meta,
	}

	return jsoniter.ConfigFastest.Marshal(data)
}

func UnmarshalCustomerNameChangedFromJSON(
	data []byte,
	streamVersion uint,
) CustomerNameChanged {

	unmarshaledData := &CustomerNameChangedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData)

	event := CustomerNameChanged{
		customerID: values.RebuildCustomerID(unmarshaledData.CustomerID),
		personName: values.RebuildPersonName(
			unmarshaledData.GivenName,
			unmarshaledData.FamilyName,
		),
		meta: EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	}

	return event
}
