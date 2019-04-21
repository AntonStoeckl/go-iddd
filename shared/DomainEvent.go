package shared

import (
    "reflect"
    "strings"
    "time"
)

type DomainEvent interface {
    Identifier() AggregateIdentifier
    EventName() string
    OccurredAt() string
}

type DomainEventMeta struct {
    Identifier AggregateIdentifier
    EventName  string
    OccurredAt string
}

func NewDomainEventMeta(id AggregateIdentifier, aggregate Aggregate, event DomainEvent) *DomainEventMeta {
    newDomainEventMeta := &DomainEventMeta{
        Identifier: id,
        EventName:  buildEventNameFor(aggregate, event),
        OccurredAt: time.Now().Format(time.RFC3339Nano),
    }

    return newDomainEventMeta
}

func buildEventNameFor(aggregate Aggregate, event DomainEvent) string {
    eventType := reflect.TypeOf(event).String()
    eventTypeParts := strings.Split(eventType, ".")
    eventName := eventTypeParts[len(eventTypeParts)-1]
    eventName = strings.Title(eventName)

    aggregateType := reflect.TypeOf(aggregate).String()
    aggregateTypeParts := strings.Split(aggregateType, ".")
    aggregateName := aggregateTypeParts[len(aggregateTypeParts)-1]
    aggregateName = strings.Title(aggregateName)

    fullEventName := aggregateName + eventName

    return fullEventName
}
