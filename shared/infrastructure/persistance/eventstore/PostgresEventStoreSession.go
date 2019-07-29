package eventstore

import (
	"database/sql"
	"go-iddd/shared"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/lib/pq"
	"golang.org/x/xerrors"
)

type PostgresEventStoreSession struct {
	tx         *sql.Tx
	eventStore *PostgresEventStore
}

func (session *PostgresEventStoreSession) LoadEventStream(
	streamID *shared.StreamID,
	fromVersion uint,
	maxEvents uint,
) (shared.DomainEvents, error) {

	queryTemplate := `SELECT event_name, payload FROM %name% 
						WHERE stream_id = $1 AND stream_version >= $2
						ORDER BY stream_version ASC
						LIMIT $3`
	query := strings.Replace(queryTemplate, "%name%", session.eventStore.tableName, 1)

	eventRows, err := session.eventStore.db.Query(query, streamID.String(), fromVersion, maxEvents)
	if err != nil {
		return nil, xerrors.Errorf(
			"postgresEventStore.LoadEventStream: %s: %w",
			err,
			shared.ErrTechnical,
		)
	}

	var stream []shared.DomainEvent
	var eventName string
	var payload string
	var domainEvent shared.DomainEvent

	for eventRows.Next() {
		if err = eventRows.Scan(&eventName, &payload); err != nil {
			return nil, xerrors.Errorf(
				"postgresEventStore.LoadEventStream: %s: %w",
				err,
				shared.ErrTechnical,
			)
		}

		if domainEvent, err = session.eventStore.unmarshal(eventName, []byte(payload)); err != nil {
			return nil, xerrors.Errorf(
				"postgresEventStore.LoadEventStream: %s: %w",
				err,
				shared.ErrUnmarshalingFailed,
			)
		}

		stream = append(stream, domainEvent)
	}

	return stream, nil
}

func (session *PostgresEventStoreSession) AppendEventsToStream(streamID *shared.StreamID, events shared.DomainEvents) error {
	queryTemplate := `INSERT INTO %name% (stream_id, stream_version, event_name, payload, occurred_at)
						VALUES ($1, $2, $3, $4, $5)`
	query := strings.Replace(queryTemplate, "%name%", session.eventStore.tableName, 1)

	for _, event := range events {
		eventJson, err := jsoniter.Marshal(event)
		if err != nil {
			return xerrors.Errorf(
				"postgresEventStoreSession.appendEventsToStreamWithTransaction: %s: %w",
				err,
				shared.ErrMarshalingFailed,
			)
		}

		_, err = session.tx.Exec(
			query,
			streamID.String(),
			event.StreamVersion(),
			event.EventName(),
			eventJson,
			event.OccurredAt(),
		)

		if err != nil {
			defaultErr := xerrors.Errorf(
				"postgresEventStoreSession.appendEventsToStreamWithTransaction: %s: %w",
				err,
				shared.ErrTechnical,
			)

			switch actualErr := err.(type) {
			case *pq.Error:
				switch actualErr.Code {
				case "23505":
					return xerrors.Errorf(
						"postgresEventStoreSession.appendEventsToStreamWithTransaction: %s: %w",
						err,
						shared.ErrConcurrencyConflict,
					)
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

func (session *PostgresEventStoreSession) Commit() error {
	if err := session.tx.Commit(); err != nil {
		return xerrors.Errorf(
			"postgresEventStoreSession.Commit: %s: %w",
			err,
			shared.ErrTechnical,
		)
	}

	return nil
}

func (session *PostgresEventStoreSession) Rollback() error {
	if err := session.tx.Rollback(); err != nil {
		return xerrors.Errorf(
			"postgresEventStoreSession.Rollback: %s: %w",
			err,
			shared.ErrTechnical,
		)
	}

	return nil
}
