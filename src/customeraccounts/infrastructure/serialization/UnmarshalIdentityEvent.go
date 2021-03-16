package serialization

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
)

// UnmarshalIdentityEvent unmarshals every know Identity event.
// It intentionally ignores unmarshaling errors, which could only happen if we would store invalid json to the EventStore.
// We have a rich test suite which would catch such issues.
func UnmarshalIdentityEvent(
	name string,
	payload []byte,
	streamVersion uint,
) (es.DomainEvent, error) {

	var event es.DomainEvent

	switch name {
	case "IdentityRegistered":
		event = unmarshalIdentityRegisteredFromJSON(payload, streamVersion)
	case "IdentityDeleted":
		event = unmarshalIdentityDeletedFromJSON(payload, streamVersion)
	default:
		err := errors.Wrapf(errors.New("event is unknown"), "unmarshalIdentityEvent [%s] failed", name)
		return nil, errors.Mark(err, shared.ErrUnmarshalingFailed)
	}

	return event, nil
}

func unmarshalIdentityRegisteredFromJSON(
	data []byte,
	streamVersion uint,
) domain.IdentityRegistered {

	unmarshaledData := &IdentityRegisteredForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := domain.RebuildIdentityRegistered(
		unmarshaledData.IdentityID,
		unmarshaledData.EmailAddress,
		unmarshaledData.ConfirmationHash,
		unmarshaledData.Password,
		es.UnmarshalEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event
}

func unmarshalIdentityDeletedFromJSON(
	data []byte,
	streamVersion uint,
) domain.IdentityDeleted {

	unmarshaledData := &IdentityDeletedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := domain.RebuildIdentityDeleted(
		unmarshaledData.IdentityID,
		es.UnmarshalEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event
}
