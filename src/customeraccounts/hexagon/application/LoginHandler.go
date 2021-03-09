package application

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
)

type LoginHandler struct {
	identities           ForStoringUniqueIdentities
	identityEventStreams ForStoringIdentityEventStreams
}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{
		identities:           nil,
		identityEventStreams: nil,
	}
}

func (h *LoginHandler) Login(emailAddress, password string) (bool, error) {
	var err error
	var identityIDValue value.IdentityID
	var emailAddressValue value.UnconfirmedEmailAddress
	var passwordValue value.PlainPassword
	var eventStream es.EventStream

	wrapWithMsg := "loginHandler.Login"

	if emailAddressValue, err = value.BuildUnconfirmedEmailAddress(emailAddress); err != nil {
		return false, errors.Wrap(err, wrapWithMsg)
	}

	if passwordValue, err = value.BuildPlainPassword(password); err != nil {
		return false, errors.Wrap(err, wrapWithMsg)
	}

	if identityIDValue, err = h.identities.FindIdentity(emailAddressValue); err != nil {
		return false, errors.Wrap(err, wrapWithMsg)
	}

	if eventStream, err = h.identityEventStreams.RetrieveIdentityEventStream(identityIDValue); err != nil {
		return false, errors.Wrap(err, wrapWithMsg)
	}

	query := domain.BuildIsMatchingPasswordForIdentity(passwordValue)

	if err = identity.IsMatchingPassword(eventStream, query); err != nil {
		return false, errors.Wrap(err, wrapWithMsg)
	}

	return true, nil
}
