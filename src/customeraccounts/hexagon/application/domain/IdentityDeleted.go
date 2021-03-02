package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type IdentityDeleted struct {
	identityID value.IdentityID
	meta       es.EventMeta
}

func BuildIdentityDeleted(
	identityID value.IdentityID,
	causationID es.MessageID,
	streamVersion uint,
) IdentityDeleted {

	event := IdentityDeleted{
		identityID: identityID,
	}

	event.meta = es.BuildEventMeta(event, causationID, streamVersion)

	return event
}

func RebuildIdentityDeleted(
	customerID string,
	meta es.EventMeta,
) IdentityDeleted {

	event := IdentityDeleted{
		identityID: value.RebuildIdentityID(customerID),
		meta:       meta,
	}

	return event
}

func (event IdentityDeleted) IdentityID() value.IdentityID {
	return event.identityID
}

func (event IdentityDeleted) Meta() es.EventMeta {
	return event.meta
}

func (event IdentityDeleted) IsFailureEvent() bool {
	return false
}

func (event IdentityDeleted) FailureReason() error {
	return nil
}
