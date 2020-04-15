package serialization

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/value"
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
) domain.CustomerRegistered {

	unmarshaledData := &CustomerRegisteredForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := domain.RebuildCustomerRegistered(
		value.RebuildCustomerID(unmarshaledData.CustomerID),
		value.RebuildEmailAddress(unmarshaledData.EmailAddress),
		value.RebuildConfirmationHash(unmarshaledData.ConfirmationHash),
		value.RebuildPersonName(
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
) domain.CustomerEmailAddressConfirmed {

	unmarshaledData := &CustomerEmailAddressConfirmedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := domain.RebuildCustomerEmailAddressConfirmed(
		value.RebuildCustomerID(unmarshaledData.CustomerID),
		value.RebuildEmailAddress(unmarshaledData.EmailAddress),
		unmarshalEventMeta(unmarshaledData.Meta, streamVersion),
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
		value.RebuildCustomerID(unmarshaledData.CustomerID),
		value.RebuildEmailAddress(unmarshaledData.EmailAddress),
		value.RebuildConfirmationHash(unmarshaledData.ConfirmationHash),
		errors.Mark(errors.New(unmarshaledData.Reason), lib.ErrDomainConstraintsViolation),
		unmarshalEventMeta(unmarshaledData.Meta, streamVersion),
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
		value.RebuildCustomerID(unmarshaledData.CustomerID),
		value.RebuildEmailAddress(unmarshaledData.EmailAddress),
		value.RebuildConfirmationHash(unmarshaledData.ConfirmationHash),
		value.RebuildEmailAddress(unmarshaledData.PreviousEmailAddress),
		unmarshalEventMeta(unmarshaledData.Meta, streamVersion),
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
		value.RebuildCustomerID(unmarshaledData.CustomerID),
		value.RebuildPersonName(
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
) domain.CustomerDeleted {

	unmarshaledData := &CustomerDeletedForJSON{}

	_ = jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData) // err intentionally ignored - see top comment

	event := domain.RebuildCustomerDeleted(
		value.RebuildCustomerID(unmarshaledData.CustomerID),
		value.RebuildEmailAddress(unmarshaledData.EmailAddress),
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
