package domain

import (
	"go-iddd/customer/domain/events"
	"go-iddd/shared"

	"golang.org/x/xerrors"
)

const unmarshalEventNamePrefix = "Customer"

func UnmarshalDomainEvent(name string, payload []byte) (shared.DomainEvent, error) {
	defaultErrFormat := "unmarshalDomainEvent [%s] failed: %w"

	switch name {
	case unmarshalEventNamePrefix + "Registered":
		event := &events.Registered{}

		if err := event.UnmarshalJSON(payload); err != nil {
			return nil, xerrors.Errorf(defaultErrFormat, name, err)
		}

		return event, nil
	case unmarshalEventNamePrefix + "EmailAddressConfirmed":
		event := &events.EmailAddressConfirmed{}

		if err := event.UnmarshalJSON(payload); err != nil {
			return nil, xerrors.Errorf(defaultErrFormat, name, err)
		}

		return event, nil
	case unmarshalEventNamePrefix + "EmailAddressChanged":
		event := &events.EmailAddressChanged{}

		if err := event.UnmarshalJSON(payload); err != nil {
			return nil, xerrors.Errorf(defaultErrFormat, name, err)
		}

		return event, nil
	default:
		return nil, xerrors.Errorf(
			"unmarshalDomainEvent [%s] failed - event is unknown: %w",
			name,
			shared.ErrUnmarshalingFailed,
		)
	}
}
