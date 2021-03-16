package es

type EventMetaForJSON struct {
	EventName   string `json:"eventName"`
	OccurredAt  string `json:"occurredAt"`
	MessageID   string `json:"messageID"`
	CausationID string `json:"causationID"`
}

func MarshalEventMeta(event DomainEvent) EventMetaForJSON {
	return EventMetaForJSON{
		EventName:   event.Meta().EventName(),
		OccurredAt:  event.Meta().OccurredAt(),
		MessageID:   event.Meta().MessageID(),
		CausationID: event.Meta().CausationID(),
	}
}
