package customer

import jsoniter "github.com/json-iterator/go"

type EmailAddressConfirmationFailed struct {
	eventName     string
	occurredAt    string
	streamVersion uint
}

func (event EmailAddressConfirmationFailed) EventName() string {
	return event.eventName
}

func (event EmailAddressConfirmationFailed) OccurredAt() string {
	return event.occurredAt
}

func (event EmailAddressConfirmationFailed) StreamVersion() uint {
	return event.streamVersion
}

func UnmarshalEmailAddressConfirmationFailedFromJSON(data []byte, streamVersion uint) EmailAddressConfirmationFailed {
	json := jsoniter.ConfigFastest

	event := EmailAddressConfirmationFailed{
		eventName:     json.Get(data, "meta").Get("eventName").ToString(),
		occurredAt:    json.Get(data, "meta").Get("occurredAt").ToString(),
		streamVersion: streamVersion,
	}

	return event
}
