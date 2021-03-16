package application

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
)

type ForStoringUniqueIdentities interface {
	FindIdentity(emailAddress value.UnconfirmedEmailAddress) (value.IdentityID, error)
	AddIdentity(identityID value.IdentityID, emailAddress value.UnconfirmedEmailAddress) error
	RemoveIdentity(identityID value.IdentityID) error
}

type ForStoringUniqueIdentitiesWithTx interface {
	ForStoringUniqueIdentities
	WithTx(tx *sql.Tx) ForStoringUniqueIdentities
}
