package eventstore

import (
	"database/sql"
	"go-iddd/shared"
	"strings"

	"golang.org/x/xerrors"
)

type PostgresEventStore struct {
	db                   *sql.DB
	tableName            string
	unmarshalDomainEvent shared.UnmarshalDomainEvent
}

func NewPostgresEventStore(
	db *sql.DB,
	tableName string,
	unmarshalDomainEvent shared.UnmarshalDomainEvent,
) *PostgresEventStore {

	return &PostgresEventStore{
		db:                   db,
		tableName:            tableName,
		unmarshalDomainEvent: unmarshalDomainEvent,
	}
}

func (store *PostgresEventStore) StartSession(tx *sql.Tx) shared.EventStore {
	return &PostgresEventStoreSession{
		tx:         tx,
		eventStore: store,
	}
}

func (store *PostgresEventStore) PurgeEventStream(streamID *shared.StreamID) error {
	queryTemplate := `DELETE FROM %name% WHERE stream_id = $1`
	query := strings.Replace(queryTemplate, "%name%", store.tableName, 1)

	_, err := store.db.Exec(query, streamID.String())

	if err != nil {
		return xerrors.Errorf(
			"postgresEventStore.PurgeEventStream: %s: %w",
			err,
			shared.ErrTechnical,
		)
	}

	return nil
}
