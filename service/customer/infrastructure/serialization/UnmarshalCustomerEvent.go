package serialization

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
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
		return nil, errors.Mark(err, lib.ErrUnmarshalingFailed)
	}

	return event, nil
}

func unmarshalCustomerRegisteredFromJSON(
	data []byte,
	streamVersion uint,
) events.CustomerRegistered {

	unmarshaledData := &CustomerRegisteredForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := events.RebuildCustomerRegistered(
		values.RebuildCustomerID(unmarshaledData.CustomerID),
		values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		values.RebuildConfirmationHash(unmarshaledData.ConfirmationHash),
		values.RebuildPersonName(
			unmarshaledData.PersonGivenName,
			unmarshaledData.PersonFamilyName,
		),
		unmarshalEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event
}

func unmarshalCustomerEmailAddressConfirmedFromJSON(
	data []byte,
	streamVersion uint,
) events.CustomerEmailAddressConfirmed {

	unmarshaledData := &CustomerEmailAddressConfirmedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := events.RebuildCustomerEmailAddressConfirmed(
		values.RebuildCustomerID(unmarshaledData.CustomerID),
		values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		unmarshalEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event
}

func unmarshalCustomerEmailAddressConfirmationFailedFromJSON(
	data []byte,
	streamVersion uint,
) events.CustomerEmailAddressConfirmationFailed {

	unmarshaledData := &CustomerEmailAddressConfirmationFailedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := events.RebuildCustomerEmailAddressConfirmationFailed(
		values.RebuildCustomerID(unmarshaledData.CustomerID),
		values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		values.RebuildConfirmationHash(unmarshaledData.ConfirmationHash),
		unmarshaledData.Reason,
		unmarshalEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event
}

func unmarshalCustomerEmailAddressChangedFromJSON(
	data []byte,
	streamVersion uint,
) events.CustomerEmailAddressChanged {

	unmarshaledData := &CustomerEmailAddressChangedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := events.RebuildCustomerEmailAddressChanged(
		values.RebuildCustomerID(unmarshaledData.CustomerID),
		values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		values.RebuildConfirmationHash(unmarshaledData.ConfirmationHash),
		values.RebuildEmailAddress(unmarshaledData.PreviousEmailAddress),
		unmarshalEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event
}

func unmarshalCustomerNameChangedFromJSON(
	data []byte,
	streamVersion uint,
) events.CustomerNameChanged {

	unmarshaledData := &CustomerNameChangedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := events.RebuildCustomerNameChanged(
		values.RebuildCustomerID(unmarshaledData.CustomerID),
		values.RebuildPersonName(
			unmarshaledData.GivenName,
			unmarshaledData.FamilyName,
		),
		unmarshalEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event
}

func unmarshalCustomerDeletedFromJSON(
	data []byte,
	streamVersion uint,
) events.CustomerDeleted {

	unmarshaledData := &CustomerDeletedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := events.RebuildCustomerDeleted(
		values.RebuildCustomerID(unmarshaledData.CustomerID),
		values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		unmarshalEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event
}

func unmarshalEventMeta(meta es.EventMetaForJSON, streamVersion uint) es.EventMeta {
	return es.RebuildEventMeta(
		meta.EventName,
		meta.OccurredAt,
		streamVersion,
	)
}
