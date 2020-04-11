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
	eventName     string
	occurredAt    string
	streamVersion uint
}

type EventMetaForJSON struct {
	EventName  string `json:"eventName"`
	OccurredAt string `json:"occurredAt"`
}

func BuildEventMeta(
	event DomainEvent,
	streamVersion uint,
) EventMeta {

	eventType := reflect.TypeOf(event).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]

	meta := EventMeta{
		eventName:     eventName,
		occurredAt:    time.Now().Format(metaTimestampFormat),
		streamVersion: streamVersion,
	}

	return meta
}

func RebuildEventMeta(
	eventName string,
	occurredAt string,
	streamVersion uint,
) EventMeta {

	return EventMeta{
		eventName:     eventName,
		occurredAt:    occurredAt,
		streamVersion: streamVersion,
	}
}

func (eventMeta EventMeta) EventName() string {
	return eventMeta.eventName
}

func (eventMeta EventMeta) OccurredAt() string {
	return eventMeta.occurredAt
}

func (eventMeta EventMeta) StreamVersion() uint {
	return eventMeta.streamVersion
}
