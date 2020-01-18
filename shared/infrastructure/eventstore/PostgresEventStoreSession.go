package eventstore

import (
	"database/sql"
	"go-iddd/shared"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/lib/pq"
)

type PostgresEventStoreSession struct {
	tx         *sql.Tx
	eventStore *PostgresEventStore
}

func (session *PostgresEventStoreSession) LoadEventStream(
	streamID shared.StreamID,
	fromVersion uint,
	maxEvents uint,
) (shared.DomainEvents, error) {

	wrapWithMsg := "postgresEventStoreSession.LoadEventStream"

	queryTemplate := `SELECT event_name, payload, stream_version FROM %name% 
						WHERE stream_id = $1 AND stream_version >= $2
						ORDER BY stream_version ASC
						LIMIT $3`
	query := strings.Replace(queryTemplate, "%name%", session.eventStore.tableName, 1)

	eventRows, err := session.eventStore.db.Query(query, streamID.String(), fromVersion, maxEvents)
	if err != nil {
		return nil, shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	var stream shared.DomainEvents
	var eventName string
	var payload string
	var streamVersion uint
	var domainEvent shared.DomainEvent

	for eventRows.Next() {
		if err = eventRows.Scan(&eventName, &payload, &streamVersion); err != nil {
			return nil, shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
		}

		if domainEvent, err = session.eventStore.unmarshalDomainEvent(eventName, []byte(payload), streamVersion); err != nil {
			return nil, shared.MarkAndWrapError(err, shared.ErrUnmarshalingFailed, wrapWithMsg)
		}

		stream = append(stream, domainEvent)
	}

	return stream, nil
}

func (session *PostgresEventStoreSession) AppendEventsToStream(
	streamID shared.StreamID,
	events shared.DomainEvents,
) error {

	wrapWithMsg := "postgresEventStoreSession.AppendEventsToStream"

	queryTemplate := `INSERT INTO %name% (stream_id, stream_version, event_name, payload, occurred_at)
						VALUES ($1, $2, $3, $4, $5)`
	query := strings.Replace(queryTemplate, "%name%", session.eventStore.tableName, 1)

	for _, event := range events {
		eventJson, err := jsoniter.Marshal(event)
		if err != nil {
			return shared.MarkAndWrapError(err, shared.ErrMarshalingFailed, wrapWithMsg)
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
			defaultErr := shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)

			switch actualErr := err.(type) {
			case *pq.Error:
				switch actualErr.Code {
				case "23505":
					return shared.MarkAndWrapError(err, shared.ErrConcurrencyConflict, wrapWithMsg)
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
