package serialization

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
)

// MarshalIdentityEvent marshals every known Identity event to json.
// It intentionally ignores marshaling errors, because they can't happen with the data types we are using.
// We have a rich test suite which would catch such issues.
func MarshalIdentityEvent(event es.DomainEvent) ([]byte, error) {
	var err error
	var json []byte

	switch actualEvent := event.(type) {
	case domain.IdentityRegistered:
		json = marshalIdentityRegistered(actualEvent)
	case domain.IdentityDeleted:
		json = marshalIdentityDeleted(actualEvent)
	default:
		err = errors.Wrapf(errors.New("event is unknown"), "marshalIdentityEvent [%s] failed", event.Meta().EventName())
		return nil, errors.Mark(err, shared.ErrMarshalingFailed)
	}

	return json, nil
}

func marshalIdentityRegistered(event domain.IdentityRegistered) []byte {
	data := IdentityRegisteredForJSON{
		IdentityID:       event.IdentityID().String(),
		EmailAddress:     event.EmailAddress().String(),
		ConfirmationHash: event.EmailAddress().ConfirmationHash().String(),
		Password:         event.Password().String(),
		Meta:             es.MarshalEventMeta(event),
	}

	json, _ := jsoniter.ConfigFastest.Marshal(data) // err intentionally ignored - see top comment

	return json
}

func marshalIdentityDeleted(event domain.IdentityDeleted) []byte {
	data := IdentityDeletedForJSON{
		IdentityID: event.IdentityID().String(),
		Meta:       es.MarshalEventMeta(event),
	}

	json, _ := jsoniter.ConfigFastest.Marshal(data) // err intentionally ignored - see top comment

	return json
}
