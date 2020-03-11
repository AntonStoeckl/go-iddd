package events

import jsoniter "github.com/json-iterator/go"

type CustomerEmailAddressChanged struct {
	emailAddress  string
	eventName     string
	occurredAt    string
	streamVersion uint
}

func (event CustomerEmailAddressChanged) EmailAddress() string {
	return event.emailAddress
}

func (event CustomerEmailAddressChanged) EventName() string {
	return event.eventName
}

func (event CustomerEmailAddressChanged) OccurredAt() string {
	return event.occurredAt
}

func (event CustomerEmailAddressChanged) StreamVersion() uint {
	return event.streamVersion
}

func UnmarshalCustomerEmailAddressChangedFromJSON(data []byte, streamVersion uint) CustomerEmailAddressChanged {
	json := jsoniter.ConfigFastest

	event := CustomerEmailAddressChanged{
		emailAddress:  json.Get(data, "emailAddress").ToString(),
		eventName:     json.Get(data, "meta").Get("eventName").ToString(),
		occurredAt:    json.Get(data, "meta").Get("occurredAt").ToString(),
		streamVersion: streamVersion,
	}

	return event
}
