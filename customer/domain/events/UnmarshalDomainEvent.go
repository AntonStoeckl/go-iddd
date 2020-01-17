package events

import (
	"go-iddd/shared"

	"github.com/cockroachdb/errors"
)

const unmarshalEventNamePrefix = "Customer"

func UnmarshalDomainEvent(name string, payload []byte, streamVersion uint) (shared.DomainEvent, error) {
	switch name {
	case unmarshalEventNamePrefix + "Registered":
		return UnmarshalRegisteredFromJSON(payload, streamVersion), nil
	case unmarshalEventNamePrefix + "EmailAddressConfirmed":
		return UnmarshalEmailAddressConfirmedFromJSON(payload, streamVersion), nil
	case unmarshalEventNamePrefix + "EmailAddressConfirmationFailed":
		return UnmarshalEmailAddressConfirmationFailedFromJSON(payload, streamVersion), nil
	case unmarshalEventNamePrefix + "EmailAddressChanged":
		return UnmarshalEmailAddressChangedFromJSON(payload, streamVersion), nil
	default:
		return nil, errors.Mark(
			errors.Wrapf(errors.New("event is unknown"), "unmarshalDomainEvent [%s] failed", name),
			shared.ErrUnmarshalingFailed,
		)
	}
}
