package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type CustomerRegistered struct {
	customerID       value.CustomerID
	emailAddress     value.UnconfirmedEmailAddress
	confirmationHash value.ConfirmationHash
	personName       value.PersonName
	meta             es.EventMeta
}

func BuildCustomerRegistered(
	customerID value.CustomerID,
	emailAddress value.UnconfirmedEmailAddress,
	confirmationHash value.ConfirmationHash,
	personName value.PersonName,
	causationID es.MessageID,
	streamVersion uint,
) CustomerRegistered {

	event := CustomerRegistered{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
		personName:       personName,
	}

	event.meta = es.BuildEventMeta(event, causationID, streamVersion)

	return event
}

func RebuildCustomerRegistered(
	customerID string,
	emailAddress string,
	confirmationHash string,
	givenName string,
	familyName string,
	meta es.EventMeta,
) CustomerRegistered {

	event := CustomerRegistered{
		customerID:       value.RebuildCustomerID(customerID),
		emailAddress:     value.RebuildUnconfirmedEmailAddress(emailAddress),
		confirmationHash: value.RebuildConfirmationHash(confirmationHash),
		personName:       value.RebuildPersonName(givenName, familyName),
		meta:             meta,
	}

	return event
}

func (event CustomerRegistered) CustomerID() value.CustomerID {
	return event.customerID
}

func (event CustomerRegistered) EmailAddress() value.UnconfirmedEmailAddress {
	return event.emailAddress
}

func (event CustomerRegistered) ConfirmationHash() value.ConfirmationHash {
	return event.confirmationHash
}

func (event CustomerRegistered) PersonName() value.PersonName {
	return event.personName
}

func (event CustomerRegistered) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerRegistered) IsFailureEvent() bool {
	return false
}

func (event CustomerRegistered) FailureReason() error {
	return nil
}
