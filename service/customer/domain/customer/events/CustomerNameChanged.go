package events

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	jsoniter "github.com/json-iterator/go"
)

type CustomerNameChanged struct {
	customerID values.CustomerID
	personName values.PersonName
	meta       EventMeta
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
	return event.meta.eventName
}

func (event CustomerNameChanged) OccurredAt() string {
	return event.meta.occurredAt
}

func (event CustomerNameChanged) StreamVersion() uint {
	return event.meta.streamVersion
}

func (event CustomerNameChanged) MarshalJSON() ([]byte, error) {
	data := struct {
		CustomerID string    `json:"customerID"`
		GivenName  string    `json:"givenName"`
		FamilyName string    `json:"familyName"`
		Meta       EventMeta `json:"meta"`
	}{
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

	anyData := jsoniter.ConfigFastest.Get(data)

	event := CustomerNameChanged{
		customerID: values.RebuildCustomerID(anyData.Get("customerID").ToString()),
		personName: values.RebuildPersonName(
			anyData.Get("givenName").ToString(),
			anyData.Get("familyName").ToString(),
		),
		meta: UnmarshalEventMetaFromJSON(data, streamVersion),
	}

	return event
}
