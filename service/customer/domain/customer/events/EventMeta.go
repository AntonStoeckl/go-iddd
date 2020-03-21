package events

import (
	"reflect"
	"strings"
	"time"

	"github.com/AntonStoeckl/go-iddd/service/lib/es"
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

func BuildEventMeta(event es.DomainEvent, streamVersion uint) EventMeta {
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

func (meta EventMeta) MarshalJSON() ([]byte, error) {
	data := &struct {
		EventName  string `json:"eventName"`
		OccurredAt string `json:"occurredAt"`
	}{
		EventName:  meta.eventName,
		OccurredAt: meta.occurredAt,
	}

	return jsoniter.ConfigFastest.Marshal(data)
}

func UnmarshalEventMetaFromJSON(
	data []byte,
	streamVersion uint,
) EventMeta {

	anyMeta := jsoniter.ConfigFastest.Get(data, "meta")

	meta := EventMeta{
		eventName:     anyMeta.Get("eventName").ToString(),
		occurredAt:    anyMeta.Get("occurredAt").ToString(),
		streamVersion: streamVersion,
	}

	return meta
}
