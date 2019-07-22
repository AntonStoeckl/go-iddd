package eventstore

import (
	"database/sql"
	"encoding/json"
	"go-iddd/shared"
	"strings"
)

type postgresEventStore struct {
	tx        *sql.Tx
	tableName string
}

func NewPostgresEventStoreSession(tx *sql.Tx, tableName string) shared.EventStore {
	return &postgresEventStore{
		tx:        tx,
		tableName: tableName,
	}
}

func (es *postgresEventStore) AppendToStream(streamID *shared.StreamID, events shared.DomainEvents) error {
	queryTemplate := `INSERT INTO eventstore (stream_id, stream_version, event_name, payload, occurred_at) VALUES ($1, $2, $3, $4, $5)`
	query := strings.Replace(queryTemplate, "%name%", es.tableName, 1)

	for _, event := range events {
		eventJson, err := json.Marshal(event)
		if err != nil {
			return err // TODO: map error
		}

		_, err = es.tx.Exec(
			query,
			streamID.String(),
			event.StreamVersion(),
			event.EventName(),
			eventJson,
			event.OccurredAt(),
		)

		if err != nil {
			return err // TODO: map error (concurrency, generic)
		}
	}

	return nil
}

func (es *postgresEventStore) LoadEventStream(streamID *shared.StreamID) (shared.DomainEvents, error) {
	panic("implement me")
}

func (es *postgresEventStore) LoadPartialEventStream(
	streamID *shared.StreamID,
	fromVersion uint,
	maxEvents uint,
) (shared.DomainEvents, error) {
	panic("implement me")
}
