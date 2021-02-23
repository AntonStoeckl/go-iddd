package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type IdentityRegistered struct {
	identityID   value.IdentityID
	emailAddress value.UnconfirmedEmailAddress
	password     value.HashedPassword
	meta         es.EventMeta
}

func BuildIdentityRegistered(
	identityID value.IdentityID,
	emailAddress value.UnconfirmedEmailAddress,
	password value.HashedPassword,
	causationID es.MessageID,
	streamVersion uint,
) IdentityRegistered {

	event := IdentityRegistered{
		identityID:   identityID,
		emailAddress: emailAddress,
		password:     password,
	}

	event.meta = es.BuildEventMeta(event, causationID, streamVersion)

	return event
}

func RebuildIdentityRegistered(
	customerID string,
	emailAddress string,
	confirmationHash string,
	meta es.EventMeta,
) IdentityRegistered {

	event := IdentityRegistered{
		identityID:   value.RebuildIdentityID(customerID),
		emailAddress: value.RebuildUnconfirmedEmailAddress(emailAddress, confirmationHash),
		meta:         meta,
	}

	return event
}

func (event IdentityRegistered) IdentityID() value.IdentityID {
	return event.identityID
}

func (event IdentityRegistered) EmailAddress() value.UnconfirmedEmailAddress {
	return event.emailAddress
}

func (event IdentityRegistered) Password() value.HashedPassword {
	return event.password
}

func (event IdentityRegistered) Meta() es.EventMeta {
	return event.meta
}

func (event IdentityRegistered) IsFailureEvent() bool {
	return false
}

func (event IdentityRegistered) FailureReason() error {
	return nil
}
