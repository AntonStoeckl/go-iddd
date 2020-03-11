package customer

import jsoniter "github.com/json-iterator/go"

type EmailAddressChanged struct {
	emailAddress  string
	eventName     string
	occurredAt    string
	streamVersion uint
}

func (event EmailAddressChanged) EmailAddress() string {
	return event.emailAddress
}

func (event EmailAddressChanged) EventName() string {
	return event.eventName
}

func (event EmailAddressChanged) OccurredAt() string {
	return event.occurredAt
}

func (event EmailAddressChanged) StreamVersion() uint {
	return event.streamVersion
}

func UnmarshalCustomerEmailAddressChangedFromJSON(data []byte, streamVersion uint) EmailAddressChanged {
	json := jsoniter.ConfigFastest

	event := EmailAddressChanged{
		emailAddress:  json.Get(data, "emailAddress").ToString(),
		eventName:     json.Get(data, "meta").Get("eventName").ToString(),
		occurredAt:    json.Get(data, "meta").Get("occurredAt").ToString(),
		streamVersion: streamVersion,
	}

	return event
}
