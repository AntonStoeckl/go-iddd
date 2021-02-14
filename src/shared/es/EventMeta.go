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
	messageID     string
	causationID   string
	streamVersion uint
}

func BuildEventMeta(
	event DomainEvent,
	causationID MessageID,
	streamVersion uint,
) EventMeta {

	meta := EventMeta{
		eventName:     buildEventName(event),
		occurredAt:    time.Now().Format(metaTimestampFormat),
		causationID:   causationID.String(),
		messageID:     GenerateMessageID().String(),
		streamVersion: streamVersion,
	}

	return meta
}

func buildEventName(event DomainEvent) string {
	eventType := reflect.TypeOf(event).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]

	return eventName
}

func RebuildEventMeta(
	eventName string,
	occurredAt string,
	messageID string,
	causationID string,
	streamVersion uint,
) EventMeta {

	return EventMeta{
		eventName:     eventName,
		occurredAt:    occurredAt,
		messageID:     messageID,
		causationID:   causationID,
		streamVersion: streamVersion,
	}
}

func (eventMeta EventMeta) EventName() string {
	return eventMeta.eventName
}

func (eventMeta EventMeta) OccurredAt() string {
	return eventMeta.occurredAt
}

func (eventMeta EventMeta) MessageID() string {
	return eventMeta.messageID
}

func (eventMeta EventMeta) CausationID() string {
	return eventMeta.causationID
}

func (eventMeta EventMeta) StreamVersion() uint {
	return eventMeta.streamVersion
}
