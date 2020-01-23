package eventstore

// import (
// 	"database/sql"
// 	"go-iddd/service/lib"
// 	"strings"
// )
//
// type PostgresEventStore struct {
// 	db                   *sql.DB
// 	tableName            string
// 	unmarshalDomainEvent lib.UnmarshalDomainEvent
// }
//
// func NewPostgresEventStore(
// 	db *sql.DB,
// 	tableName string,
// 	unmarshalDomainEvent lib.UnmarshalDomainEvent,
// ) *PostgresEventStore {
//
// 	return &PostgresEventStore{
// 		db:                   db,
// 		tableName:            tableName,
// 		unmarshalDomainEvent: unmarshalDomainEvent,
// 	}
// }
//
// func (store *PostgresEventStore) StartSession(tx *sql.Tx) lib.EventStore {
// 	return &PostgresEventStoreSession{
// 		tx:         tx,
// 		eventStore: store,
// 	}
// }
//
// func (store *PostgresEventStore) PurgeEventStream(streamID lib.StreamID) error {
// 	queryTemplate := `DELETE FROM %name% WHERE stream_id = $1`
// 	query := strings.Replace(queryTemplate, "%name%", store.tableName, 1)
//
// 	_, err := store.db.Exec(query, streamID.String())
//
// 	if err != nil {
// 		return lib.MarkAndWrapError(err, lib.ErrTechnical, "postgresEventStore.PurgeEventStream")
// 	}
//
// 	return nil
// }
