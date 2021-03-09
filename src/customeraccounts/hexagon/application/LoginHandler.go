package application

import "github.com/cockroachdb/errors"

type LoginHandler struct {
	identities           ForStoringUniqueIdentities
	identityEventStreams ForStoringIdentityEventStreams
}

func NewLoginHandler(identityEventStreams ForStoringIdentityEventStreams) *LoginHandler {
	return &LoginHandler{
		identities:           nil,
		identityEventStreams: identityEventStreams,
	}
}

func (h *LoginHandler) Login(emailAddress, password string) (bool, error) {
	//var err error
	//var identityIDValue value.IdentityID
	//var emailAddressValue value.UnconfirmedEmailAddress
	//var passwordValue value.PlainPassword
	//var eventStream es.EventStream
	//
	//wrapWithMsg := "loginHandler.Login"
	//
	//if emailAddressValue, err = value.BuildUnconfirmedEmailAddress(emailAddress); err != nil {
	//	return false, errors.Wrap(err, wrapWithMsg)
	//}
	//
	//if passwordValue, err = value.BuildPlainPassword(password); err != nil {
	//	return false, errors.Wrap(err, wrapWithMsg)
	//}
	//
	//if identityIDValue, err = h.identities.FindIdentity(emailAddressValue); err != nil {
	//	return false, errors.Wrap(err, wrapWithMsg)
	//}
	//
	//if eventStream, err = h.identityEventStreams.RetrieveIdentityEventStream(identityIDValue); err != nil {
	//	return false, errors.Wrap(err, wrapWithMsg)
	//}
	//
	//query := domain.BuildIsMatchingPasswordForIdentity(passwordValue)
	//
	//if err = identity.IsMatchingPassword(eventStream, query); err != nil {
	//	return false, errors.Wrap(err, wrapWithMsg)
	//}

	return false, errors.New("dummy error")
}
