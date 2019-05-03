package shared

import (
	"errors"
	"reflect"
	"strings"
	"time"
)

type DomainEvent interface {
	Identifier() string
	EventName() string
	OccurredAt() string
}

func NewDomainEventMeta(aggregateID string, event DomainEvent, aggregateName string) *DomainEventMeta {
	newDomainEventMeta := &DomainEventMeta{
		Identifier: aggregateID,
		EventName:  buildEventNameFor(event, aggregateName),
		OccurredAt: time.Now().Format(time.RFC3339Nano),
	}

	return newDomainEventMeta
}

func buildEventNameFor(event DomainEvent, withAggregateName string) string {
	eventType := reflect.TypeOf(event).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]
	eventName = strings.Title(eventName)

	fullEventName := withAggregateName + eventName

	return fullEventName
}

type DomainEventMeta struct {
	Identifier string `json:"identifier"`
	EventName  string `json:"eventName"`
	OccurredAt string `json:"occurredAt"`
}

func UnmarshalDomainEventMeta(data interface{}) (*DomainEventMeta, error) {
	values, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("zefix")
	}

	meta := &DomainEventMeta{}

	for key, value := range values {
		value, ok := value.(string)
		if !ok {
			return nil, errors.New("zefix")
		}

		switch key {
		case "identifier":
			meta.Identifier = value
		case "eventName":
			meta.EventName = value
		case "occurredAt":
			meta.OccurredAt = value
		}
	}

	return meta, nil
}
