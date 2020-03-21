package events

import (
	"reflect"
	"strings"
	"time"

	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

const (
	metaTimestampFormat = time.RFC3339Nano
)

type EventMeta struct {
	EventName     string `json:"eventName"`
	OccurredAt    string `json:"occurredAt"`
	streamVersion uint
}

func BuildEventMeta(
	event es.DomainEvent,
	streamVersion uint,
) EventMeta {

	eventType := reflect.TypeOf(event).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]

	meta := EventMeta{
		EventName:     eventName,
		OccurredAt:    time.Now().Format(metaTimestampFormat),
		streamVersion: streamVersion,
	}

	return meta
}

func EnrichEventMeta(
	eventMeta EventMeta,
	streamVersion uint,
) EventMeta {

	eventMeta.streamVersion = streamVersion

	return eventMeta
}
