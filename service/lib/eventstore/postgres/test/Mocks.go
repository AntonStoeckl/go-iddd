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

func (event SomeEvent) EventName() string {
	return event.name
}

func (event SomeEvent) OccurredAt() string {
	return event.occurredAt
}

func (event SomeEvent) IndicatesAnError() (bool, string) {
	return false, ""
}

func (event SomeEvent) StreamVersion() uint {
	return event.version
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

func (event BrokenMarshalingEvent) EventName() string {
	return event.name
}

func (event BrokenMarshalingEvent) OccurredAt() string {
	return event.occurredAt
}

func (event BrokenMarshalingEvent) IndicatesAnError() (bool, string) {
	return false, ""
}

func (event BrokenMarshalingEvent) StreamVersion() uint {
	return event.version
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

func (event BrokenUnmarshalingEvent) EventName() string {
	return event.name
}

func (event BrokenUnmarshalingEvent) OccurredAt() string {
	return event.occurredAt
}

func (event BrokenUnmarshalingEvent) IndicatesAnError() (bool, string) {
	return false, ""
}

func (event BrokenUnmarshalingEvent) StreamVersion() uint {
	return event.version
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
			Name:       actualEvent.EventName(),
			OccurredAt: actualEvent.OccurredAt(),
			Version:    actualEvent.StreamVersion(),
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
			Name:       actualEvent.EventName(),
			OccurredAt: actualEvent.OccurredAt(),
			Version:    actualEvent.StreamVersion(),
		}

		return jsoniter.Marshal(data)
	default:
		return nil, errors.New("mocked marshaling error unknown event")
	}
}
