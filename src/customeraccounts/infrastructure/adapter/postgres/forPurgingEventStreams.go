package postgres

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type forPurgingEventStreams func(streamID es.StreamID, tx *sql.Tx) error
