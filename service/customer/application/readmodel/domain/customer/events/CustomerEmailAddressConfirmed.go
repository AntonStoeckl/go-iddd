package events

import jsoniter "github.com/json-iterator/go"

type CustomerEmailAddressConfirmed struct {
	eventName     string
	occurredAt    string
	streamVersion uint
}

func (event CustomerEmailAddressConfirmed) EventName() string {
	return event.eventName
}

func (event CustomerEmailAddressConfirmed) OccurredAt() string {
	return event.occurredAt
}

func (event CustomerEmailAddressConfirmed) StreamVersion() uint {
	return event.streamVersion
}

func UnmarshalCustomerEmailAddressConfirmedFromJSON(data []byte, streamVersion uint) CustomerEmailAddressConfirmed {
	json := jsoniter.ConfigFastest

	event := CustomerEmailAddressConfirmed{
		eventName:     json.Get(data, "meta").Get("eventName").ToString(),
		occurredAt:    json.Get(data, "meta").Get("occurredAt").ToString(),
		streamVersion: streamVersion,
	}

	return event
}
