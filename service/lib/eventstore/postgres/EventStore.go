package postgres

import (
	"database/sql"
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/lib/pq"
)

type EventStore struct {
	db                   *sql.DB
	tableName            string
	unmarshalDomainEvent es.UnmarshalDomainEvent
}

func NewEventStore(
	db *sql.DB,
	tableName string,
	unmarshalDomainEvent es.UnmarshalDomainEvent,
) *EventStore {
	store := &EventStore{
		db:                   db,
		tableName:            tableName,
		unmarshalDomainEvent: unmarshalDomainEvent,
	}

	return store
}

func (eventStore *EventStore) AppendEventsToStream(
	streamID es.StreamID,
	events es.DomainEvents,
	tx *sql.Tx,
) error {

	wrapWithMsg := "eventStore.AppendEventsToStream"

	queryTemplate := `INSERT INTO %name% (stream_id, stream_version, event_name, payload, occurred_at)
						VALUES ($1, $2, $3, $4, $5)`
	query := strings.Replace(queryTemplate, "%name%", eventStore.tableName, 1)

	for _, event := range events {
		eventJson, err := jsoniter.Marshal(event)
		if err != nil {
			return lib.MarkAndWrapError(err, lib.ErrMarshalingFailed, wrapWithMsg)
		}

		_, err = tx.Exec(
			query,
			streamID.String(),
			event.StreamVersion(),
			event.EventName(),
			eventJson,
			event.OccurredAt(),
		)

		if err != nil {
			defaultErr := lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)

			switch actualErr := err.(type) {
			case *pq.Error:
				switch actualErr.Code {
				case "23505":
					return lib.MarkAndWrapError(err, lib.ErrConcurrencyConflict, wrapWithMsg)
				default:
					return defaultErr // some other postgres error (e.g. table does not exist)
				}
			default:
				return defaultErr // some other DB error (e.g. tx already closed, no connection)
			}
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
