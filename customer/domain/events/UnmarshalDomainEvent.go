package events

import (
	"go-iddd/shared"

	"github.com/cockroachdb/errors"
)

const unmarshalEventNamePrefix = "Customer"

func UnmarshalDomainEvent(name string, payload []byte) (shared.DomainEvent, error) {
	defaultErrFormat := "unmarshalDomainEvent [%s] failed"

	switch name {
	case unmarshalEventNamePrefix + "Registered":
		event := &Registered{}

		if err := event.UnmarshalJSON(payload); err != nil {
			return nil, errors.Wrapf(err, defaultErrFormat, name)
		}

		return event, nil
	case unmarshalEventNamePrefix + "EmailAddressConfirmed":
		event := &EmailAddressConfirmed{}

		if err := event.UnmarshalJSON(payload); err != nil {
			return nil, errors.Wrapf(err, defaultErrFormat, name)
		}

		return event, nil
	case unmarshalEventNamePrefix + "EmailAddressConfirmationFailed":
		event := &EmailAddressConfirmationFailed{}

		if err := event.UnmarshalJSON(payload); err != nil {
			return nil, errors.Wrapf(err, defaultErrFormat, name)
		}

		return event, nil
	case unmarshalEventNamePrefix + "EmailAddressChanged":
		event := &EmailAddressChanged{}

		if err := event.UnmarshalJSON(payload); err != nil {
			return nil, errors.Wrapf(err, defaultErrFormat, name)
		}

		return event, nil
	default:
		return nil, errors.Mark(
			errors.Wrapf(errors.New("event is unknown"), defaultErrFormat, name),
			shared.ErrUnmarshalingFailed,
		)
	}
}
