package hexagon

import "github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"

type ForRegisteringIdentities func(identityIDValue value.IdentityID, emailAddress string, plainPassword string) error
