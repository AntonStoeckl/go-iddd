package application

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity"
	"github.com/cockroachdb/errors"
)

type LoginHandler struct {
	uniqueIdentities     ForStoringUniqueIdentities
	identityEventStreams ForStoringIdentityEventStreams
}

func NewLoginHandler(
	uniqueIdentities ForStoringUniqueIdentities,
	identityEventStreams ForStoringIdentityEventStreams,
) *LoginHandler {

	return &LoginHandler{
		uniqueIdentities:     uniqueIdentities,
		identityEventStreams: identityEventStreams,
	}
}

func (h *LoginHandler) HandleLogIn(
	emailAddress string,
	password string,
) (bool, error) {

	logIn := func() error {
		query, err := domain.BuildIsMatchingPasswordForIdentity(emailAddress, password)
		if err != nil {
			return err
		}

		identityIDValue, err := h.uniqueIdentities.FindIdentity(query.EmailAddress())
		if err != nil {
			return err
		}

		eventStream, err := h.identityEventStreams.RetrieveEventStream(identityIDValue)
		if err != nil {
			return err
		}

		if err := identity.IsMatchingPassword(eventStream, query); err != nil {
			return err
		}

		return nil
	}

	if err := logIn(); err != nil {
		return false, errors.Wrap(err, "LoginHandler.HandleLogIn")
	}

	return true, nil
}
