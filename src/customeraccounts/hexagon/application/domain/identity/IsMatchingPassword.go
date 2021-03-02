package identity

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
)

func IsMatchingPassword(stream es.EventStream, query domain.IsMatchingPasswordForIdentity) error {
	identity := buildCurrentStateFrom(stream)

	if err := assertNotDeleted(identity); err != nil {
		return errors.Wrap(shared.ErrInvalidCredentials, "isMatchingPassword")
	}

	if !identity.password.CompareWith(query.SuppliedPassword()) {
		return errors.Wrap(shared.ErrInvalidCredentials, "isMatchingPassword")
	}

	return nil
}
