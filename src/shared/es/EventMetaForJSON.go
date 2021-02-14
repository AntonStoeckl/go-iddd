package es

type EventMetaForJSON struct {
	EventName   string `json:"eventName"`
	OccurredAt  string `json:"occurredAt"`
	MessageID   string `json:"messageID"`
	CausationID string `json:"causationID"`
}
