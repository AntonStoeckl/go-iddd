package events

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
	jsoniter "github.com/json-iterator/go"
)

func UnmarshalCustomerEvent(
	name string,
	payload []byte,
	streamVersion uint,
) (es.DomainEvent, error) {

	var err error
	var event es.DomainEvent

	switch name {
	case "CustomerRegistered":
		event, err = unmarshalCustomerRegisteredFromJSON(payload, streamVersion)
	case "CustomerEmailAddressConfirmed":
		event, err = unmarshalCustomerEmailAddressConfirmedFromJSON(payload, streamVersion)
	case "CustomerEmailAddressConfirmationFailed":
		event, err = unmarshalCustomerEmailAddressConfirmationFailedFromJSON(payload, streamVersion)
	case "CustomerEmailAddressChanged":
		event, err = unmarshalCustomerEmailAddressChangedFromJSON(payload, streamVersion)
	case "CustomerNameChanged":
		event, err = unmarshalCustomerNameChangedFromJSON(payload, streamVersion)
	case "CustomerDeleted":
		event, err = unmarshalCustomerDeletedFromJSON(payload, streamVersion)
	default:
		err = errors.Wrapf(errors.New("event is unknown"), "unmarshalCustomerEvent [%s] failed", name)
	}

	if err != nil {
		return nil, errors.Mark(err, lib.ErrUnmarshalingFailed)
	}

	return event, nil
}

func unmarshalCustomerRegisteredFromJSON(
	data []byte,
	streamVersion uint,
) (CustomerRegistered, error) {

	unmarshaledData := &CustomerRegisteredForJSON{}

	if err := jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData); err != nil {
		return CustomerRegistered{}, err
	}

	event := RebuildCustomerRegistered(
		values.RebuildCustomerID(unmarshaledData.CustomerID),
		values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		values.RebuildConfirmationHash(unmarshaledData.ConfirmationHash),
		values.RebuildPersonName(
			unmarshaledData.PersonGivenName,
			unmarshaledData.PersonFamilyName,
		),
		es.EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event, nil
}

func unmarshalCustomerEmailAddressConfirmedFromJSON(
	data []byte,
	streamVersion uint,
) (CustomerEmailAddressConfirmed, error) {

	unmarshaledData := &CustomerEmailAddressConfirmedForJSON{}

	if err := jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData); err != nil {
		return CustomerEmailAddressConfirmed{}, err
	}

	event := RebuildCustomerEmailAddressConfirmed(
		values.RebuildCustomerID(unmarshaledData.CustomerID),
		values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		es.EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event, nil
}

func unmarshalCustomerEmailAddressConfirmationFailedFromJSON(
	data []byte,
	streamVersion uint,
) (CustomerEmailAddressConfirmationFailed, error) {

	unmarshaledData := &CustomerEmailAddressConfirmationFailedForJSON{}

	if err := jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData); err != nil {
		return CustomerEmailAddressConfirmationFailed{}, err
	}

	event := RebuildCustomerEmailAddressConfirmationFailed(
		values.RebuildCustomerID(unmarshaledData.CustomerID),
		values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		values.RebuildConfirmationHash(unmarshaledData.ConfirmationHash),
		"", // TODO
		es.EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event, nil
}

func unmarshalCustomerEmailAddressChangedFromJSON(
	data []byte,
	streamVersion uint,
) (CustomerEmailAddressChanged, error) {

	unmarshaledData := &CustomerEmailAddressChangedForJSON{}

	if err := jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData); err != nil {
		return CustomerEmailAddressChanged{}, err
	}

	event := RebuildCustomerEmailAddressChanged(
		values.RebuildCustomerID(unmarshaledData.CustomerID),
		values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		values.RebuildConfirmationHash(unmarshaledData.ConfirmationHash),
		values.RebuildEmailAddress(unmarshaledData.PreviousEmailAddress),
		es.EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event, nil
}

func unmarshalCustomerNameChangedFromJSON(
	data []byte,
	streamVersion uint,
) (CustomerNameChanged, error) {

	unmarshaledData := &CustomerNameChangedForJSON{}

	if err := jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData); err != nil {
		return CustomerNameChanged{}, err
	}

	event := RebuildCustomerNameChanged(
		values.RebuildCustomerID(unmarshaledData.CustomerID),
		values.RebuildPersonName(
			unmarshaledData.GivenName,
			unmarshaledData.FamilyName,
		),
		es.EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event, nil
}

func unmarshalCustomerDeletedFromJSON(
	data []byte,
	streamVersion uint,
) (CustomerDeleted, error) {

	unmarshaledData := &CustomerDeletedForJSON{}

	if err := jsoniter.ConfigFastest.Unmarshal(data, unmarshaledData); err != nil {
		return CustomerDeleted{}, err
	}

	event := RebuildCustomerDeleted(
		values.RebuildCustomerID(unmarshaledData.CustomerID),
		values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		es.EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	)

	return event, nil
}
