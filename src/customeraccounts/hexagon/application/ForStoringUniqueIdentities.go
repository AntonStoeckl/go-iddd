package application

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
)

type ForStoringUniqueIdentities interface {
	FindIdentity(emailAddress value.UnconfirmedEmailAddress) (value.IdentityID, error)
}
