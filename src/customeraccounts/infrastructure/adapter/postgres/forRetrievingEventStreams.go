package postgres

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type forRetrievingEventStreams func(streamID es.StreamID, fromVersion uint, maxEvents uint, db *sql.DB) (es.EventStream, error)
