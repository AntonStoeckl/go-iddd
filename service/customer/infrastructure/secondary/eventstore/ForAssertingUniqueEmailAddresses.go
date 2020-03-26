package eventstore

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type ForAssertingUniqueEmailAddresses interface {
	Assert(recordedEvents es.DomainEvents, tx *sql.Tx) error
}
