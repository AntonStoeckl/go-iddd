package events

import (
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
)

func MarshalCustomerEvent(event es.DomainEvent) ([]byte, error) {
	var err error
	var json []byte

	switch actualEvent := event.(type) {
	case CustomerRegistered:
		json, err = marshalCustomerRegistered(actualEvent)
	case CustomerEmailAddressConfirmed:
		json, err = marshalCustomerEmailAddressConfirmed(actualEvent)
	case CustomerEmailAddressConfirmationFailed:
		json, err = marshalCustomerEmailAddressConfirmationFailed(actualEvent)
	case CustomerEmailAddressChanged:
		json, err = marshalCustomerEmailAddressChanged(actualEvent)
	case CustomerNameChanged:
		json, err = marshalCustomerNameChanged(actualEvent)
	case CustomerDeleted:
		json, err = marshalCustomerDeleted(actualEvent)
	default:
		err = errors.Wrapf(errors.New("event is unknown"), "marshalCustomerEvent [%s] failed", event.EventName())
	}

	if err != nil {
		return nil, errors.Mark(err, lib.ErrMarshalingFailed)
	}

	return json, nil
}

func marshalCustomerRegistered(event CustomerRegistered) ([]byte, error) {
	data := CustomerRegisteredForJSON{
		CustomerID:       event.CustomerID().String(),
		EmailAddress:     event.EmailAddress().String(),
		ConfirmationHash: event.ConfirmationHash().String(),
		PersonGivenName:  event.PersonName().GivenName(),
		PersonFamilyName: event.PersonName().FamilyName(),
		Meta:             event.Meta(),
	}

	return jsoniter.ConfigFastest.Marshal(data)
}

func marshalCustomerEmailAddressConfirmed(event CustomerEmailAddressConfirmed) ([]byte, error) {
	data := CustomerEmailAddressConfirmedForJSON{
		CustomerID:   event.CustomerID().String(),
		EmailAddress: event.EmailAddress().String(),
		Meta:         event.Meta(),
	}

	return jsoniter.ConfigFastest.Marshal(data)
}

func marshalCustomerEmailAddressConfirmationFailed(event CustomerEmailAddressConfirmationFailed) ([]byte, error) {
	data := CustomerEmailAddressConfirmationFailedForJSON{
		CustomerID:       event.CustomerID().String(),
		EmailAddress:     event.EmailAddress().String(),
		ConfirmationHash: event.ConfirmationHash().String(),
		Meta:             event.Meta(),
	}

	return jsoniter.ConfigFastest.Marshal(data)
}

func marshalCustomerEmailAddressChanged(event CustomerEmailAddressChanged) ([]byte, error) {
	data := CustomerEmailAddressChangedForJSON{
		CustomerID:           event.CustomerID().String(),
		EmailAddress:         event.EmailAddress().String(),
		ConfirmationHash:     event.ConfirmationHash().String(),
		PreviousEmailAddress: event.PreviousEmailAddress().String(),
		Meta:                 event.Meta(),
	}

	return jsoniter.ConfigFastest.Marshal(data)
}

func marshalCustomerNameChanged(event CustomerNameChanged) ([]byte, error) {
	data := CustomerNameChangedForJSON{
		CustomerID: event.CustomerID().String(),
		GivenName:  event.PersonName().GivenName(),
		FamilyName: event.PersonName().FamilyName(),
		Meta:       event.Meta(),
	}

	return jsoniter.ConfigFastest.Marshal(data)
}

func marshalCustomerDeleted(event CustomerDeleted) ([]byte, error) {
	data := CustomerDeletedForJSON{
		CustomerID:   event.CustomerID().String(),
		EmailAddress: event.EmailAddress().String(),
		Meta:         event.Meta(),
	}

	return jsoniter.ConfigFastest.Marshal(data)
}
