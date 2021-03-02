package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type IsMatchingPasswordForIdentity struct {
	suppliedPassword value.PlainPassword
	messageID        es.MessageID
}

func BuildIsMatchingPasswordForIdentity(suppliedPassword value.PlainPassword) IsMatchingPasswordForIdentity {
	return IsMatchingPasswordForIdentity{
		suppliedPassword: suppliedPassword,
		messageID:        es.GenerateMessageID(),
	}
}

func (query IsMatchingPasswordForIdentity) SuppliedPassword() value.PlainPassword {
	return query.suppliedPassword
}

func (query IsMatchingPasswordForIdentity) MessageID() es.MessageID {
	return query.messageID
}
