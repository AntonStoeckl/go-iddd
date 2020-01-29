package lib

import (
	"database/sql"
)

type EventStore interface {
	AppendEventsToStream(streamID StreamID, events DomainEvents, tx *sql.Tx) error
	LoadEventStream(streamID StreamID, fromVersion uint, maxEvents uint) (DomainEvents, error)
	PurgeEventStream(streamID StreamID) error
}
