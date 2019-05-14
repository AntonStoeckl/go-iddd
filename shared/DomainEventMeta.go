package shared

import (
	"reflect"
	"strings"
	"time"
)

const DomainEventMetaTimestampFormat = time.RFC3339Nano

type DomainEventMeta struct {
	Identifier string `json:"identifier"`
	EventName  string `json:"eventName"`
	OccurredAt string `json:"occurredAt"`
}

func NewDomainEventMeta(aggregateID string, event DomainEvent, aggregateName string) *DomainEventMeta {
	eventType := reflect.TypeOf(event).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]
	eventName = strings.Title(eventName)
	fullEventName := aggregateName + eventName

	newDomainEventMeta := &DomainEventMeta{
		Identifier: aggregateID,
		EventName:  fullEventName,
		OccurredAt: time.Now().Format(DomainEventMetaTimestampFormat),
	}

	return newDomainEventMeta
}
