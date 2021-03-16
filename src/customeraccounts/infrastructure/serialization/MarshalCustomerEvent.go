package serialization

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
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
	case domain.CustomerRegistered:
		json = marshalCustomerRegistered(actualEvent)
	case domain.CustomerEmailAddressConfirmed:
		json = marshalCustomerEmailAddressConfirmed(actualEvent)
	case domain.CustomerEmailAddressConfirmationFailed:
		json = marshalCustomerEmailAddressConfirmationFailed(actualEvent)
	case domain.CustomerEmailAddressChanged:
		json = marshalCustomerEmailAddressChanged(actualEvent)
	case domain.CustomerNameChanged:
		json = marshalCustomerNameChanged(actualEvent)
	case domain.CustomerDeleted:
		json = marshalCustomerDeleted(actualEvent)
	default:
		err = errors.Wrapf(errors.New("event is unknown"), "marshalCustomerEvent [%s] failed", event.Meta().EventName())
		return nil, errors.Mark(err, shared.ErrMarshalingFailed)
	}

	return json, nil
}

func marshalCustomerRegistered(event domain.CustomerRegistered) []byte {
	data := CustomerRegisteredForJSON{
		CustomerID:       event.CustomerID().String(),
		EmailAddress:     event.EmailAddress().String(),
		ConfirmationHash: event.EmailAddress().ConfirmationHash().String(),
		PersonGivenName:  event.PersonName().GivenName(),
		PersonFamilyName: event.PersonName().FamilyName(),
		Meta:             es.MarshalEventMeta(event),
	}

	json, _ := jsoniter.ConfigFastest.Marshal(data) // err intentionally ignored - see top comment

	return json
}

func marshalCustomerEmailAddressConfirmed(event domain.CustomerEmailAddressConfirmed) []byte {
	data := CustomerEmailAddressConfirmedForJSON{
		CustomerID:   event.CustomerID().String(),
		EmailAddress: event.EmailAddress().String(),
		Meta:         es.MarshalEventMeta(event),
	}

	json, _ := jsoniter.ConfigFastest.Marshal(data) // err intentionally ignored - see top comment

	return json
}

func marshalCustomerEmailAddressConfirmationFailed(event domain.CustomerEmailAddressConfirmationFailed) []byte {
	data := CustomerEmailAddressConfirmationFailedForJSON{
		CustomerID:       event.CustomerID().String(),
		ConfirmationHash: event.ConfirmationHash().String(),
		Reason:           event.FailureReason().Error(),
		Meta:             es.MarshalEventMeta(event),
	}

	json, _ := jsoniter.ConfigFastest.Marshal(data) // err intentionally ignored - see top comment

	return json
}

func marshalCustomerEmailAddressChanged(event domain.CustomerEmailAddressChanged) []byte {
	data := CustomerEmailAddressChangedForJSON{
		CustomerID:       event.CustomerID().String(),
		EmailAddress:     event.EmailAddress().String(),
		ConfirmationHash: event.EmailAddress().ConfirmationHash().String(),
		Meta:             es.MarshalEventMeta(event),
	}

	json, _ := jsoniter.ConfigFastest.Marshal(data) // err intentionally ignored - see top comment

	return json
}

func marshalCustomerNameChanged(event domain.CustomerNameChanged) []byte {
	data := CustomerNameChangedForJSON{
		CustomerID: event.CustomerID().String(),
		GivenName:  event.PersonName().GivenName(),
		FamilyName: event.PersonName().FamilyName(),
		Meta:       es.MarshalEventMeta(event),
	}

	json, _ := jsoniter.ConfigFastest.Marshal(data) // err intentionally ignored - see top comment

	return json
}

func marshalCustomerDeleted(event domain.CustomerDeleted) []byte {
	data := CustomerDeletedForJSON{
		CustomerID: event.CustomerID().String(),
		Meta:       es.MarshalEventMeta(event),
	}

	json, _ := jsoniter.ConfigFastest.Marshal(data) // err intentionally ignored - see top comment

	return json
}
