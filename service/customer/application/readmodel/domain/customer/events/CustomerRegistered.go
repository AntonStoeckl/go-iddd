package events

import (
	jsoniter "github.com/json-iterator/go"
)

type CustomerRegistered struct {
	customerID    string
	emailAddress  string
	givenName     string
	familyName    string
	eventName     string
	occurredAt    string
	streamVersion uint
}

func (event CustomerRegistered) CustomerID() string {
	return event.customerID
}

func (event CustomerRegistered) EmailAddress() string {
	return event.emailAddress
}

func (event CustomerRegistered) GivenName() string {
	return event.givenName
}

func (event CustomerRegistered) FamilyName() string {
	return event.familyName
}

func (event CustomerRegistered) EventName() string {
	return event.eventName
}

func (event CustomerRegistered) OccurredAt() string {
	return event.occurredAt
}

func (event CustomerRegistered) StreamVersion() uint {
	return event.streamVersion
}

func UnmarshalCustomerRegisteredFromJSON(data []byte, streamVersion uint) CustomerRegistered {
	json := jsoniter.ConfigFastest

	event := CustomerRegistered{
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
