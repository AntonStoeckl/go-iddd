package customer

import (
	jsoniter "github.com/json-iterator/go"
)

type Registered struct {
	customerID    string
	emailAddress  string
	givenName     string
	familyName    string
	eventName     string
	occurredAt    string
	streamVersion uint
}

func (event Registered) CustomerID() string {
	return event.customerID
}

func (event Registered) EmailAddress() string {
	return event.emailAddress
}

func (event Registered) GivenName() string {
	return event.givenName
}

func (event Registered) FamilyName() string {
	return event.familyName
}

func (event Registered) EventName() string {
	return event.eventName
}

func (event Registered) OccurredAt() string {
	return event.occurredAt
}

func (event Registered) StreamVersion() uint {
	return event.streamVersion
}

func UnmarshalCustomerRegisteredFromJSON(data []byte, streamVersion uint) Registered {
	json := jsoniter.ConfigFastest

	event := Registered{
		customerID:    json.Get(data, "customerID").ToString(),
		emailAddress:  json.Get(data, "emailAddress").ToString(),
		givenName:     json.Get(data, "personGivenName").ToString(),
		familyName:    json.Get(data, "personFamilyName").ToString(),
		eventName:     json.Get(data, "meta").Get("eventName").ToString(),
		occurredAt:    json.Get(data, "meta").Get("occurredAt").ToString(),
		streamVersion: streamVersion,
	}

	return event
}
