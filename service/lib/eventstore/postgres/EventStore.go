package postgres

import (
	"database/sql"
	"strings"

	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
	"github.com/lib/pq"
)

type EventStore struct {
	db                   *sql.DB
	tableName            string
	marshalDomainEvent   es.MarshalDomainEvent
	unmarshalDomainEvent es.UnmarshalDomainEvent
}

func NewEventStore(
	db *sql.DB,
	tableName string,
	marshalDomainEvent es.MarshalDomainEvent,
	unmarshalDomainEvent es.UnmarshalDomainEvent,
) *EventStore {
	store := &EventStore{
		db:                   db,
		tableName:            tableName,
		marshalDomainEvent:   marshalDomainEvent,
		unmarshalDomainEvent: unmarshalDomainEvent,
	}

	return store
}

func (eventStore *EventStore) AppendEventsToStream(
	streamID es.StreamID,
	events es.DomainEvents,
	tx *sql.Tx,
) error {

	var err error

	wrapWithMsg := "eventStore.AppendEventsToStream"

	queryTemplate := `INSERT INTO %name% (stream_id, stream_version, event_name, occurred_at, payload)
						VALUES ($1, $2, $3, $4, $5)`
	query := strings.Replace(queryTemplate, "%name%", eventStore.tableName, 1)

	for _, event := range events {
		var eventJson []byte

		eventJson, err = eventStore.marshalDomainEvent(event)
		if err != nil {
			return lib.MarkAndWrapError(err, lib.ErrMarshalingFailed, wrapWithMsg)
		}

		_, err = tx.Exec(
			query,
			streamID.String(),
			event.Meta().StreamVersion(),
			event.Meta().EventName(),
			event.Meta().OccurredAt(),
			eventJson,
		)

		if err != nil {
			return errors.Wrap(eventStore.mapPostgresErrors(err), wrapWithMsg)
		}
	}

	return nil
}

func (eventStore *EventStore) LoadEventStream(
	streamID es.StreamID,
	fromVersion uint,
	maxEvents uint,
) (es.DomainEvents, error) {

	wrapWithMsg := "eventStore.LoadEventStream"

	queryTemplate := `SELECT event_name, payload, stream_version FROM %name% 
						WHERE stream_id = $1 AND stream_version >= $2
						ORDER BY stream_version ASC
						LIMIT $3`

	query := strings.Replace(queryTemplate, "%name%", eventStore.tableName, 1)

	eventRows, err := eventStore.db.Query(query, streamID.String(), fromVersion, maxEvents)
	if err != nil {
		return nil, lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	var eventStream es.DomainEvents
	var eventName string
	var payload string
	var streamVersion uint
	var domainEvent es.DomainEvent

	for eventRows.Next() {
		if eventRows.Err() != nil {
			return nil, lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
		}

		if err = eventRows.Scan(&eventName, &payload, &streamVersion); err != nil {
			return nil, lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
		}

		if domainEvent, err = eventStore.unmarshalDomainEvent(eventName, []byte(payload), streamVersion); err != nil {
			return nil, lib.MarkAndWrapError(err, lib.ErrUnmarshalingFailed, wrapWithMsg)
		}

		eventStream = append(eventStream, domainEvent)
	}

	return eventStream, nil
}

func (eventStore *EventStore) PurgeEventStream(streamID es.StreamID) error {
	queryTemplate := `DELETE FROM %name% WHERE stream_id = $1`
	query := strings.Replace(queryTemplate, "%name%", eventStore.tableName, 1)

	_, err := eventStore.db.Exec(query, streamID.String())

	if err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, "eventStore.PurgeEventStream")
	}

	return nil
}

func (eventStore *EventStore) mapPostgresErrors(err error) error {
	defaultErr := errors.Mark(err, lib.ErrTechnical)

	switch actualErr := err.(type) {
	case *pq.Error:
		switch actualErr.Code {
		case "23505":
			return errors.Mark(err, lib.ErrConcurrencyConflict)
		default:
			return defaultErr // some other postgres error (e.g. table does not exist)
		}
	default:
		return defaultErr // some other DB error (e.g. tx already closed, no connection)
	}
}
