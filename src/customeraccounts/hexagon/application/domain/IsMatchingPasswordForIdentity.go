package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type IsMatchingPasswordForIdentity struct {
	emailAddress value.UnconfirmedEmailAddress
	password     value.PlainPassword
	messageID    es.MessageID
}

func BuildIsMatchingPasswordForIdentity(
	emailAddress string,
	password string,
) (IsMatchingPasswordForIdentity, error) {

	emailAddressValue, err := value.BuildUnconfirmedEmailAddress(emailAddress)
	if err != nil {
		return IsMatchingPasswordForIdentity{}, err
	}

	passwordValue, err := value.BuildPlainPassword(password)
	if err != nil {
		return IsMatchingPasswordForIdentity{}, err
	}

	return IsMatchingPasswordForIdentity{
		emailAddress: emailAddressValue,
		password:     passwordValue,
		messageID:    es.GenerateMessageID(),
	}, nil
}

func (query IsMatchingPasswordForIdentity) EmailAddress() value.UnconfirmedEmailAddress {
	return query.emailAddress
}

func (query IsMatchingPasswordForIdentity) SuppliedPassword() value.PlainPassword {
	return query.password
}

func (query IsMatchingPasswordForIdentity) MessageID() es.MessageID {
	return query.messageID
}
