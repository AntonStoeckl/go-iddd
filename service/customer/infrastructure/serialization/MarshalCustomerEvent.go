package serialization

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
)

// MarshalCustomerEvent marshals every known Customer event to json.
// It intentionally ignores marshaling errors, because they can't happen with the data types we are using.
// We have a rich test suite which would catch such issues.
func MarshalCustomerEvent(event es.DomainEvent) ([]byte, error) {
	var err error
	var json []byte

	switch actualEvent := event.(type) {
	case events.CustomerRegistered:
		json = marshalCustomerRegistered(actualEvent)
	case events.CustomerEmailAddressConfirmed:
		json = marshalCustomerEmailAddressConfirmed(actualEvent)
	case events.CustomerEmailAddressConfirmationFailed:
		json = marshalCustomerEmailAddressConfirmationFailed(actualEvent)
	case events.CustomerEmailAddressChanged:
		json = marshalCustomerEmailAddressChanged(actualEvent)
	case events.CustomerNameChanged:
		json = marshalCustomerNameChanged(actualEvent)
	case events.CustomerDeleted:
		json = marshalCustomerDeleted(actualEvent)
	default:
		err = errors.Wrapf(errors.New("event is unknown"), "marshalCustomerEvent [%s] failed", event.EventName())
		return nil, errors.Mark(err, lib.ErrMarshalingFailed)
	}

	return json, nil
}

func marshalCustomerRegistered(event events.CustomerRegistered) []byte {
	data := CustomerRegisteredForJSON{
		CustomerID:       event.CustomerID().String(),
		EmailAddress:     event.EmailAddress().String(),
		ConfirmationHash: event.ConfirmationHash().String(),
		PersonGivenName:  event.PersonName().GivenName(),
		PersonFamilyName: event.PersonName().FamilyName(),
		Meta:             event.Meta(),
	}

	json, _ := jsoniter.ConfigFastest.Marshal(data) // err intentionally ignored - see top comment

	return json
}

func marshalCustomerEmailAddressConfirmed(event events.CustomerEmailAddressConfirmed) []byte {
	data := CustomerEmailAddressConfirmedForJSON{
		CustomerID:   event.CustomerID().String(),
		EmailAddress: event.EmailAddress().String(),
		Meta:         event.Meta(),
	}

	json, _ := jsoniter.ConfigFastest.Marshal(data) // err intentionally ignored - see top comment

	return json
}

func marshalCustomerEmailAddressConfirmationFailed(event events.CustomerEmailAddressConfirmationFailed) []byte {
	_, reason := event.IndicatesAnError()

	data := CustomerEmailAddressConfirmationFailedForJSON{
		CustomerID:       event.CustomerID().String(),
		EmailAddress:     event.EmailAddress().String(),
		ConfirmationHash: event.ConfirmationHash().String(),
		Reason:           reason,
		Meta:             event.Meta(),
	}

	json, _ := jsoniter.ConfigFastest.Marshal(data) // err intentionally ignored - see top comment

	return json
}

func marshalCustomerEmailAddressChanged(event events.CustomerEmailAddressChanged) []byte {
	data := CustomerEmailAddressChangedForJSON{
		CustomerID:           event.CustomerID().String(),
		EmailAddress:         event.EmailAddress().String(),
		ConfirmationHash:     event.ConfirmationHash().String(),
		PreviousEmailAddress: event.PreviousEmailAddress().String(),
		Meta:                 event.Meta(),
	}

	json, _ := jsoniter.ConfigFastest.Marshal(data) // err intentionally ignored - see top comment

	return json
}

func marshalCustomerNameChanged(event events.CustomerNameChanged) []byte {
	data := CustomerNameChangedForJSON{
		CustomerID: event.CustomerID().String(),
		GivenName:  event.PersonName().GivenName(),
		FamilyName: event.PersonName().FamilyName(),
		Meta:       event.Meta(),
	}

	json, _ := jsoniter.ConfigFastest.Marshal(data) // err intentionally ignored - see top comment

	return json
}

func marshalCustomerDeleted(event events.CustomerDeleted) []byte {
	data := CustomerDeletedForJSON{
		CustomerID:   event.CustomerID().String(),
		EmailAddress: event.EmailAddress().String(),
		Meta:         event.Meta(),
	}

	json, _ := jsoniter.ConfigFastest.Marshal(data) // err intentionally ignored - see top comment

	return json
}
