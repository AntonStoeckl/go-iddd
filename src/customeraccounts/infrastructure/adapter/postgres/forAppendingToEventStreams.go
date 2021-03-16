package postgres

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type forAppendingEventsToStreams func(
	streamID es.StreamID,
	events []es.DomainEvent,
	marshalDomainEvent es.MarshalDomainEvent,
	tx *sql.Tx,
) error
