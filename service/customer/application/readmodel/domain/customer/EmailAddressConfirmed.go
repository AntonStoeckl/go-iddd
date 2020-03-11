package customer

import jsoniter "github.com/json-iterator/go"

type EmailAddressConfirmed struct {
	eventName     string
	occurredAt    string
	streamVersion uint
}

func (event EmailAddressConfirmed) EventName() string {
	return event.eventName
}

func (event EmailAddressConfirmed) OccurredAt() string {
	return event.occurredAt
}

func (event EmailAddressConfirmed) StreamVersion() uint {
	return event.streamVersion
}

func UnmarshalCustomerEmailAddressConfirmedFromJSON(data []byte, streamVersion uint) EmailAddressConfirmed {
	json := jsoniter.ConfigFastest

	event := EmailAddressConfirmed{
		eventName:     json.Get(data, "meta").Get("eventName").ToString(),
		occurredAt:    json.Get(data, "meta").Get("occurredAt").ToString(),
		streamVersion: streamVersion,
	}

	return event
}
