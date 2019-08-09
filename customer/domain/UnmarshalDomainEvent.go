package domain

import (
	"go-iddd/customer/domain/events"
	"go-iddd/shared"

	"github.com/cockroachdb/errors"
)

const unmarshalEventNamePrefix = "Customer"

func UnmarshalDomainEvent(name string, payload []byte) (shared.DomainEvent, error) {
	defaultErrFormat := "unmarshalDomainEvent [%s] failed"

	switch name {
	case unmarshalEventNamePrefix + "Registered":
		event := &events.Registered{}

		if err := event.UnmarshalJSON(payload); err != nil {
			return nil, errors.Wrapf(err, defaultErrFormat, name)
		}

		return event, nil
	case unmarshalEventNamePrefix + "EmailAddressConfirmed":
		event := &events.EmailAddressConfirmed{}

		if err := event.UnmarshalJSON(payload); err != nil {
			return nil, errors.Wrapf(err, defaultErrFormat, name)
		}

		return event, nil
	case unmarshalEventNamePrefix + "EmailAddressChanged":
		event := &events.EmailAddressChanged{}

		if err := event.UnmarshalJSON(payload); err != nil {
			return nil, errors.Wrapf(err, defaultErrFormat, name)
		}

		return event, nil
	default:
		return nil, errors.Mark(
			errors.Newf("unmarshalDomainEvent [%s] failed - event is unknown", name),
			shared.ErrUnmarshalingFailed,
		)
	}
}
