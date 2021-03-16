package serialization

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
)

// UnmarshalCustomerEvent unmarshals every know Customer event.
// It intentionally ignores unmarshaling errors, which could only happen if we would store invalid json to the EventStore.
// We have a rich test suite which would catch such issues.
func UnmarshalCustomerEvent(
	name string,
	payload []byte,
	streamVersion uint,
) (es.DomainEvent, error) {

	var event es.DomainEvent

	switch name {
	case "CustomerRegistered":
		event = unmarshalCustomerRegisteredFromJSON(payload, streamVersion)
	case "CustomerEmailAddressConfirmed":
		event = unmarshalCustomerEmailAddressConfirmedFromJSON(payload, streamVersion)
	case "CustomerEmailAddressConfirmationFailed":
		event = unmarshalCustomerEmailAddressConfirmationFailedFromJSON(payload, streamVersion)
	case "CustomerEmailAddressChanged":
		event = unmarshalCustomerEmailAddressChangedFromJSON(payload, streamVersion)
	case "CustomerNameChanged":
		event = unmarshalCustomerNameChangedFromJSON(payload, streamVersion)
	case "CustomerDeleted":
		event = unmarshalCustomerDeletedFromJSON(payload, streamVersion)
	default:
		err := errors.Wrapf(errors.New("event is unknown"), "unmarshalCustomerEvent [%s] failed", name)
		return nil, errors.Mark(err, shared.ErrUnmarshalingFailed)
	}

	return event, nil
}

func unmarshalCustomerRegisteredFromJSON(
	data []byte,
	streamVersion uint,
) domain.CustomerRegistered {

	unmarshaledData := &CustomerRegisteredForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := domain.RebuildCustomerRegistered(
		unmarshaledData.CustomerID,
		unmarshaledData.EmailAddress,
		unmarshaledData.ConfirmationHash,
		unmarshaledData.PersonGivenName,
		unmarshaledData.PersonFamilyName,
		es.UnmarshalEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event
}

func unmarshalCustomerEmailAddressConfirmedFromJSON(
	data []byte,
	streamVersion uint,
) domain.CustomerEmailAddressConfirmed {

	unmarshaledData := &CustomerEmailAddressConfirmedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := domain.RebuildCustomerEmailAddressConfirmed(
		unmarshaledData.CustomerID,
		unmarshaledData.EmailAddress,
		es.UnmarshalEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event
}

func unmarshalCustomerEmailAddressConfirmationFailedFromJSON(
	data []byte,
	streamVersion uint,
) domain.CustomerEmailAddressConfirmationFailed {

	unmarshaledData := &CustomerEmailAddressConfirmationFailedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := domain.RebuildCustomerEmailAddressConfirmationFailed(
		unmarshaledData.CustomerID,
		unmarshaledData.ConfirmationHash,
		unmarshaledData.Reason,
		es.UnmarshalEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event
}

func unmarshalCustomerEmailAddressChangedFromJSON(
	data []byte,
	streamVersion uint,
) domain.CustomerEmailAddressChanged {

	unmarshaledData := &CustomerEmailAddressChangedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := domain.RebuildCustomerEmailAddressChanged(
		unmarshaledData.CustomerID,
		unmarshaledData.EmailAddress,
		unmarshaledData.ConfirmationHash,
		es.UnmarshalEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event
}

func unmarshalCustomerNameChangedFromJSON(
	data []byte,
	streamVersion uint,
) domain.CustomerNameChanged {

	unmarshaledData := &CustomerNameChangedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := domain.RebuildCustomerNameChanged(
		unmarshaledData.CustomerID,
		unmarshaledData.GivenName,
		unmarshaledData.FamilyName,
		es.UnmarshalEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event
}

func unmarshalCustomerDeletedFromJSON(
	data []byte,
	streamVersion uint,
) domain.CustomerDeleted {

	unmarshaledData := &CustomerDeletedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := domain.RebuildCustomerDeleted(
		unmarshaledData.CustomerID,
		es.UnmarshalEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event
}
