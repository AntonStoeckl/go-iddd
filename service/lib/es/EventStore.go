package es

import "database/sql"

type EventStore interface {
	AppendEventsToStream(streamID StreamID, events RecordedEvents, tx *sql.Tx) error
	LoadEventStream(streamID StreamID, fromVersion uint, maxEvents uint) (EventStream, error)
	PurgeEventStream(streamID StreamID) error
}
