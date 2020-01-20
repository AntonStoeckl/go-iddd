package eventstore

import (
	"database/sql"
	"go-iddd/service/lib"
)

type StartsEventStoreSessions interface {
	StartSession(tx *sql.Tx) lib.EventStore
}
