package es

import (
	"reflect"
	"strings"
	"time"
)

const (
	metaTimestampFormat = time.RFC3339Nano
)

type EventMeta struct {
	EventName     string `json:"eventName"`
	OccurredAt    string `json:"occurredAt"`
	StreamVersion uint   `json:"-"`
}

func BuildEventMeta(
	event DomainEvent,
	streamVersion uint,
) EventMeta {

	eventType := reflect.TypeOf(event).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]

	meta := EventMeta{
		EventName:     eventName,
		OccurredAt:    time.Now().Format(metaTimestampFormat),
		StreamVersion: streamVersion,
	}

	return meta
}

func EnrichEventMeta(
	eventMeta EventMeta,
	streamVersion uint,
) EventMeta {

	eventMeta.StreamVersion = streamVersion

	return eventMeta
}
