package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type RegisterIdentity struct {
	identityID   value.IdentityID
	emailAddress value.UnconfirmedEmailAddress
	password     value.HashedPassword
	messageID    es.MessageID
}

func BuildRegisterIdentity(
	identityID value.IdentityID,
	emailAddress string,
	password string,
) (RegisterIdentity, error) {

	emailAddressValue, err := value.BuildUnconfirmedEmailAddress(emailAddress)
	if err != nil {
		return RegisterIdentity{}, err
	}

	plainPasswordValue, err := value.BuildPlainPassword(password)
	if err != nil {
		return RegisterIdentity{}, err
	}

	hashedPasswordValue, err := value.HashedPasswordFromPlainPassword(plainPasswordValue)
	if err != nil {
		return RegisterIdentity{}, err
	}

	command := RegisterIdentity{
		identityID:   identityID,
		emailAddress: emailAddressValue,
		password:     hashedPasswordValue,
		messageID:    es.GenerateMessageID(),
	}

	return command, nil
}

func (command RegisterIdentity) IdentityID() value.IdentityID {
	return command.identityID
}

func (command RegisterIdentity) EmailAddress() value.UnconfirmedEmailAddress {
	return command.emailAddress
}

func (command RegisterIdentity) Password() value.HashedPassword {
	return command.password
}

func (command RegisterIdentity) MessageID() es.MessageID {
	return command.messageID
}
