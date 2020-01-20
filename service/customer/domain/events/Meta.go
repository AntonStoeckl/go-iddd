package events

import (
	jsoniter "github.com/json-iterator/go"
)

type Meta struct {
	identifier    string
	eventName     string
	occurredAt    string
	streamVersion uint
}

func (meta Meta) MarshalJSON() ([]byte, error) {
	data := &struct {
		Identifier string `json:"identifier"`
		EventName  string `json:"eventName"`
		OccurredAt string `json:"occurredAt"`
	}{
		Identifier: meta.identifier,
		EventName:  meta.eventName,
		OccurredAt: meta.occurredAt,
	}

	return jsoniter.Marshal(data)
}

func UnmarshalMetaFromJSON(data []byte, streamVersion uint) Meta {
	anyMeta := jsoniter.Get(data, "meta")

	meta := Meta{
		identifier:    anyMeta.Get("identifier").ToString(),
		eventName:     anyMeta.Get("eventName").ToString(),
		occurredAt:    anyMeta.Get("occurredAt").ToString(),
		streamVersion: streamVersion,
	}

	return meta
}
