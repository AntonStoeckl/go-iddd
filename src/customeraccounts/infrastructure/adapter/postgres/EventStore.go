package postgres

import (
	"database/sql"
	"strings"

	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	"github.com/lib/pq"
)

type EventStore struct {
	eventStoreTableName  string
	marshalDomainEvent   es.MarshalDomainEvent
	unmarshalDomainEvent es.UnmarshalDomainEvent
}

func NewEventStore(
	eventStoreTableName string,
	marshalDomainEvent es.MarshalDomainEvent,
	unmarshalDomainEvent es.UnmarshalDomainEvent,
) *EventStore {

	return &EventStore{
		eventStoreTableName:  eventStoreTableName,
		marshalDomainEvent:   marshalDomainEvent,
		unmarshalDomainEvent: unmarshalDomainEvent,
	}
}

func (s *EventStore) RetrieveEventStream(
	streamID es.StreamID,
	fromVersion uint,
	maxEvents uint,
	db *sql.DB,
) (es.EventStream, error) {

	var err error
	wrapWithMsg := "retrieveEventStream"

	queryTemplate := `SELECT event_name, payload, stream_version FROM %name% 
						WHERE stream_id = $1 AND stream_version >= $2
						ORDER BY stream_version ASC
						LIMIT $3`

	query := strings.Replace(queryTemplate, "%name%", s.eventStoreTableName, 1)

	eventRows, err := db.Query(query, streamID.String(), fromVersion, maxEvents)
	if err != nil {
		return nil, shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	defer eventRows.Close()

	var eventStream es.EventStream
	var eventName string
	var payload string
	var streamVersion uint
	var domainEvent es.DomainEvent

	for eventRows.Next() {
		if eventRows.Err() != nil {
			return nil, shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
		}

		if err = eventRows.Scan(&eventName, &payload, &streamVersion); err != nil {
			return nil, shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
		}

		if domainEvent, err = s.unmarshalDomainEvent(eventName, []byte(payload), streamVersion); err != nil {
			return nil, shared.MarkAndWrapError(err, shared.ErrUnmarshalingFailed, wrapWithMsg)
		}

		eventStream = append(eventStream, domainEvent)
	}

	return eventStream, nil
}

func (s *EventStore) AppendEventsToStream(
	streamID es.StreamID,
	events []es.DomainEvent,
	tx *sql.Tx,
) error {

	var err error
	wrapWithMsg := "appendEventsToStream"

	queryTemplate := `INSERT INTO %name% (stream_id, stream_version, event_name, occurred_at, payload)
						VALUES ($1, $2, $3, $4, $5)`
	query := strings.Replace(queryTemplate, "%name%", s.eventStoreTableName, 1)

	for _, event := range events {
		var eventJSON []byte

		eventJSON, err = s.marshalDomainEvent(event)
		if err != nil {
			return shared.MarkAndWrapError(err, shared.ErrMarshalingFailed, wrapWithMsg)
		}

		_, err = tx.Exec(
			query,
			streamID.String(),
			event.Meta().StreamVersion(),
			event.Meta().EventName(),
			event.Meta().OccurredAt(),
			eventJSON,
		)

		if err != nil {
			return errors.Wrap(s.mapEventStorePostgresErrors(err), wrapWithMsg)
		}
	}

	return nil
}

func (s *EventStore) PurgeEventStream(
	streamID es.StreamID,
	tx *sql.Tx,
) error {

	queryTemplate := `DELETE FROM %name% WHERE stream_id = $1`
	query := strings.Replace(queryTemplate, "%name%", s.eventStoreTableName, 1)

	if _, err := tx.Exec(query, streamID.String()); err != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, "purgeEventStream")
	}

	return nil
}

func (s *EventStore) mapEventStorePostgresErrors(err error) error {
	// nolint:errorlint // errors.As() suggested, but somehow cockroachdb/errors can't convert this properly
	if actualErr, ok := err.(*pq.Error); ok {
		if actualErr.Code == "23505" {
			return errors.Mark(err, shared.ErrConcurrencyConflict)
		}
	}

	return errors.Mark(err, shared.ErrTechnical) // some other DB error (Tx closed, wrong table, ...)
}
