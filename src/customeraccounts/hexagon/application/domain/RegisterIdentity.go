package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type RegisterIdentity struct {
	identityID   value.IdentityID
	emailAddress value.UnconfirmedEmailAddress
	messageID    es.MessageID
}

func BuildRegisterIdentity(
	identityID value.IdentityID,
	emailAddress value.UnconfirmedEmailAddress,
) RegisterIdentity {

	command := RegisterIdentity{
		identityID:   identityID,
		emailAddress: emailAddress,
		messageID:    es.GenerateMessageID(),
	}

	return command
}

func (command RegisterIdentity) IdentityID() value.IdentityID {
	return command.identityID
}

func (command RegisterIdentity) EmailAddress() value.UnconfirmedEmailAddress {
	return command.emailAddress
}

func (command RegisterIdentity) MessageID() es.MessageID {
	return command.messageID
}
