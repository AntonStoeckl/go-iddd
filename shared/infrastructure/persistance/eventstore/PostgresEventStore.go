package eventstore

import (
	"database/sql"
	"go-iddd/shared"
	"strings"

	"golang.org/x/xerrors"
)

type PostgresEventStore struct {
	db        *sql.DB
	tableName string
	unmarshal shared.UnmarshalDomainEvent
}

func NewPostgresEventStore(
	db *sql.DB,
	tableName string,
	unmarshal shared.UnmarshalDomainEvent,
) *PostgresEventStore {

	return &PostgresEventStore{
		db:        db,
		tableName: tableName,
		unmarshal: unmarshal,
	}
}

func (store *PostgresEventStore) StartSession() (shared.EventStoreSession, error) {
	tx, err := store.db.Begin()
	if err != nil {
		return nil, xerrors.Errorf(
			"postgresEventStore.StartSession: %s: %w",
			err,
			shared.ErrTechnical,
		)
	}

	session := &PostgresEventStoreSession{
		tx:         tx,
		eventStore: store,
	}

	return session, nil
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
