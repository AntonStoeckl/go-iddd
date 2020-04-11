// +build test

package test

import (
	"time"

	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
)

/*** mocked CustomerID ***/

type SomeID struct {
	Value string
}

func (someID SomeID) ID() string {
	return someID.Value
}

/*** mocked Event that works ***/

type SomeEvent struct {
	id         SomeID
	name       string
	version    uint
	occurredAt string
}

func CreateSomeEvent(forId SomeID, withVersion uint) SomeEvent {
	return SomeEvent{
		id:         forId,
		name:       "SomeEvent",
		version:    withVersion,
		occurredAt: time.Now().Format(time.RFC3339Nano),
	}
}

func (event SomeEvent) Meta() es.EventMeta {
	return es.RebuildEventMeta(event.name, event.occurredAt, event.version)
}

func (event SomeEvent) IsFailureEvent() bool {
	return false
}

func (event SomeEvent) FailureReason() error {
	return nil
}

func UnmarshalSomeEventFromJSON(data []byte) SomeEvent {
	someEvent := SomeEvent{
		id:         SomeID{Value: jsoniter.Get(data, "customerID").ToString()},
		name:       jsoniter.Get(data, "name").ToString(),
		version:    jsoniter.Get(data, "version").ToUint(),
		occurredAt: jsoniter.Get(data, "occurredAt").ToString(),
	}

	return someEvent
}

/*** mocked Event with broken marshaling ***/

type BrokenMarshalingEvent struct {
	id         SomeID
	name       string
	version    uint
	occurredAt string
}

func CreateBrokenMarshalingEvent(forId SomeID, withVersion uint) BrokenMarshalingEvent {
	return BrokenMarshalingEvent{
		id:         forId,
		name:       "BrokenMarshalingEvent",
		version:    withVersion,
		occurredAt: time.Now().Format(time.RFC3339Nano),
	}
}

func (event BrokenMarshalingEvent) Meta() es.EventMeta {
	return es.RebuildEventMeta(event.name, event.occurredAt, event.version)
}

func (event BrokenMarshalingEvent) IsFailureEvent() bool {
	return false
}

func (event BrokenMarshalingEvent) FailureReason() error {
	return nil
}

/*** mocked Event with broken unmarshaling ***/

type BrokenUnmarshalingEvent struct {
	id         SomeID
	name       string
	version    uint
	occurredAt string
}

func CreateBrokenUnmarshalingEvent(forId SomeID, withVersion uint) BrokenUnmarshalingEvent {
	return BrokenUnmarshalingEvent{
		id:         forId,
		name:       "BrokenUnmarshalingEvent",
		version:    withVersion,
		occurredAt: time.Now().Format(time.RFC3339Nano),
	}
}

func (event BrokenUnmarshalingEvent) Meta() es.EventMeta {
	return es.RebuildEventMeta(event.name, event.occurredAt, event.version)
}

func (event BrokenUnmarshalingEvent) IsFailureEvent() bool {
	return false
}

func (event BrokenUnmarshalingEvent) FailureReason() error {
	return nil
}

/*** Unmarshal mocked events ***/

func UnmarshalMockEvents(name string, payload []byte, streamVersion uint) (es.DomainEvent, error) {
	_ = streamVersion
	defaultErrFormat := "unmarshalDomainEvent [%s] failed: %w"

	switch name {
	case "SomeEvent":
		return UnmarshalSomeEventFromJSON(payload), nil
	case "BrokenUnmarshalingEvent":
		return nil, errors.Errorf(defaultErrFormat, name, errors.New("mocked marshaling error"))
	default:
		return nil, errors.New("unknown mocked event to unmarshal")
	}
}

func MarshalMockEvents(event es.DomainEvent) ([]byte, error) {
	switch actualEvent := event.(type) {
	case SomeEvent:
		data := &struct {
			ID         string `json:"customerID"`
			Name       string `json:"name"`
			Version    uint   `json:"version"`
			OccurredAt string `json:"occurredAt"`
		}{
			ID:         actualEvent.id.Value,
			Name:       actualEvent.Meta().EventName(),
			OccurredAt: actualEvent.Meta().OccurredAt(),
			Version:    actualEvent.Meta().StreamVersion(),
		}

		return jsoniter.Marshal(data)
	case BrokenMarshalingEvent:
		return nil, errors.New("mocked marshaling error")
	case BrokenUnmarshalingEvent:
		data := &struct {
			ID         string `json:"customerID"`
			Name       string `json:"name"`
			Version    uint   `json:"version"`
			OccurredAt string `json:"occurredAt"`
		}{
			ID:         actualEvent.id.Value,
			Name:       actualEvent.Meta().EventName(),
			OccurredAt: actualEvent.Meta().OccurredAt(),
			Version:    actualEvent.Meta().StreamVersion(),
		}

		return jsoniter.Marshal(data)
	default:
		return nil, errors.New("mocked marshaling error unknown event")
	}
}
