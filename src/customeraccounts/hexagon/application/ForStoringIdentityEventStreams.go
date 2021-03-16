package application

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type ForStoringIdentityEventStreams interface {
	RetrieveEventStream(id value.IdentityID) (es.EventStream, error)
	StartEventStream(identityRegistered domain.IdentityRegistered) error
	PurgeEventStream(id value.IdentityID) error
}

type ForStoringIdentityEventStreamsWithTx interface {
	ForStoringIdentityEventStreams
	WithTx(tx *sql.Tx) ForStoringIdentityEventStreams
}
