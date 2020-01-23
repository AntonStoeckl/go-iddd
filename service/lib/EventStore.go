package lib

import (
	"database/sql"
)

type EventStore interface {
	LoadEventStream(streamID StreamID, fromVersion uint, maxEvents uint) (DomainEvents, error)
	AppendEventsToStream(streamID StreamID, events DomainEvents, tx *sql.Tx) error
	PurgeEventStream(streamID StreamID) error
}
