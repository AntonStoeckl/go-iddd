package eventstore

import (
	"database/sql"
	"go-iddd/service/lib"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/lib/pq"
)

type PostgresEventStore struct {
	db                   *sql.DB
	tableName            string
	unmarshalDomainEvent lib.UnmarshalDomainEvent
}

func NewPostgresEventStore(
	db *sql.DB,
	tableName string,
	unmarshalDomainEvent lib.UnmarshalDomainEvent,
) *PostgresEventStore {
	store := &PostgresEventStore{
		db:                   db,
		tableName:            tableName,
		unmarshalDomainEvent: unmarshalDomainEvent,
	}

	return store
}

func (es *PostgresEventStore) AppendEventsToStream(
	streamID lib.StreamID,
	events lib.DomainEvents,
	tx *sql.Tx,
) error {

	wrapWithMsg := "postgresEventStoreSession.AppendEventsToStream"

	queryTemplate := `INSERT INTO %name% (stream_id, stream_version, event_name, payload, occurred_at)
						VALUES ($1, $2, $3, $4, $5)`
	query := strings.Replace(queryTemplate, "%name%", es.tableName, 1)

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

func (es *PostgresEventStore) LoadEventStream(
	streamID lib.StreamID,
	fromVersion uint,
	maxEvents uint,
) (lib.DomainEvents, error) {

	wrapWithMsg := "postgresEventStoreSession.LoadEventStream"

	queryTemplate := `SELECT event_name, payload, stream_version FROM %name% 
						WHERE stream_id = $1 AND stream_version >= $2
						ORDER BY stream_version ASC
						LIMIT $3`

	query := strings.Replace(queryTemplate, "%name%", es.tableName, 1)

	eventRows, err := es.db.Query(query, streamID.String(), fromVersion, maxEvents)
	if err != nil {
		return nil, lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	var stream lib.DomainEvents
	var eventName string
	var payload string
	var streamVersion uint
	var domainEvent lib.DomainEvent

	for eventRows.Next() {
		if err = eventRows.Scan(&eventName, &payload, &streamVersion); err != nil {
			return nil, lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
		}

		if domainEvent, err = es.unmarshalDomainEvent(eventName, []byte(payload), streamVersion); err != nil {
			return nil, lib.MarkAndWrapError(err, lib.ErrUnmarshalingFailed, wrapWithMsg)
		}

		stream = append(stream, domainEvent)
	}

	return stream, nil
}

func (es *PostgresEventStore) PurgeEventStream(streamID lib.StreamID) error {
	queryTemplate := `DELETE FROM %name% WHERE stream_id = $1`
	query := strings.Replace(queryTemplate, "%name%", es.tableName, 1)

	_, err := es.db.Exec(query, streamID.String())

	if err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, "postgresEventStore.PurgeEventStream")
	}

	return nil
}
