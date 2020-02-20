package events

import (
	"go-iddd/service/lib/es"
	"reflect"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

const (
	metaTimestampFormat = time.RFC3339Nano
)

type EventMeta struct {
	eventName     string
	occurredAt    string
	streamVersion uint
}

func BuildEventMeta(event es.DomainEvent, prefix string, streamVersion uint) EventMeta {
	eventType := reflect.TypeOf(event).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]
	fullEventName := prefix + eventName

	meta := EventMeta{
		eventName:     fullEventName,
		occurredAt:    time.Now().Format(metaTimestampFormat),
		streamVersion: streamVersion,
	}

	return meta
}

func (meta EventMeta) MarshalJSON() ([]byte, error) {
	data := &struct {
		EventName  string `json:"eventName"`
		OccurredAt string `json:"occurredAt"`
	}{
		EventName:  meta.eventName,
		OccurredAt: meta.occurredAt,
	}

	return jsoniter.Marshal(data)
}

func UnmarshalEventMetaFromJSON(data []byte, streamVersion uint) EventMeta {
	anyMeta := jsoniter.Get(data, "meta")

	meta := EventMeta{
		eventName:     anyMeta.Get("eventName").ToString(),
		occurredAt:    anyMeta.Get("occurredAt").ToString(),
		streamVersion: streamVersion,
	}

	return meta
}
