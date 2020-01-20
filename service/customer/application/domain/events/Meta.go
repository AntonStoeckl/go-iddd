package events

import (
	jsoniter "github.com/json-iterator/go"
)

type Meta struct {
	eventName     string
	occurredAt    string
	streamVersion uint
}

func (meta Meta) MarshalJSON() ([]byte, error) {
	data := &struct {
		EventName  string `json:"eventName"`
		OccurredAt string `json:"occurredAt"`
	}{
		EventName:  meta.eventName,
		OccurredAt: meta.occurredAt,
	}

	return jsoniter.Marshal(data)
}

func UnmarshalMetaFromJSON(data []byte, streamVersion uint) Meta {
	anyMeta := jsoniter.Get(data, "meta")

	meta := Meta{
		eventName:     anyMeta.Get("eventName").ToString(),
		occurredAt:    anyMeta.Get("occurredAt").ToString(),
		streamVersion: streamVersion,
	}

	return meta
}
