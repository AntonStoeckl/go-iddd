package eventstore

import (
	"database/sql"
	"go-iddd/shared"
)

type StartsEventStoreSessions interface {
	StartSession(tx *sql.Tx) shared.EventStore
}
