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

	event := CustomerRegistered{
		customerID:       values.RebuildCustomerID(unmarshaledData.CustomerID),
		emailAddress:     values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		confirmationHash: values.RebuildConfirmationHash(unmarshaledData.ConfirmationHash),
		personName: values.RebuildPersonName(
			unmarshaledData.PersonGivenName,
			unmarshaledData.PersonFamilyName,
		),
		meta: es.EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	}

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

	event := CustomerEmailAddressConfirmed{
		customerID:   values.RebuildCustomerID(unmarshaledData.CustomerID),
		emailAddress: values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		meta:         es.EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	}

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

	event := CustomerEmailAddressConfirmationFailed{
		customerID:       values.RebuildCustomerID(unmarshaledData.CustomerID),
		emailAddress:     values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		confirmationHash: values.RebuildConfirmationHash(unmarshaledData.ConfirmationHash),
		meta:             es.EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	}

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

	event := CustomerEmailAddressChanged{
		customerID:           values.RebuildCustomerID(unmarshaledData.CustomerID),
		emailAddress:         values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		confirmationHash:     values.RebuildConfirmationHash(unmarshaledData.ConfirmationHash),
		previousEmailAddress: values.RebuildEmailAddress(unmarshaledData.PreviousEmailAddress),
		meta:                 es.EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	}

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

	event := CustomerNameChanged{
		customerID: values.RebuildCustomerID(unmarshaledData.CustomerID),
		personName: values.RebuildPersonName(
			unmarshaledData.GivenName,
			unmarshaledData.FamilyName,
		),
		meta: es.EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	}

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

	event := CustomerDeleted{
		customerID:   values.RebuildCustomerID(unmarshaledData.CustomerID),
		emailAddress: values.RebuildEmailAddress(unmarshaledData.EmailAddress),
		meta:         es.EnrichEventMeta(unmarshaledData.Meta, streamVersion),
	}

	return event, nil
}
