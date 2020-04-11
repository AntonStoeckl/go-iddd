package events

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

type CustomerRegistered struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	personName       values.PersonName
	meta             es.EventMeta
}

func BuildCustomerRegistered(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	confirmationHash values.ConfirmationHash,
	personName values.PersonName,
	streamVersion uint,
) CustomerRegistered {

	event := CustomerRegistered{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
		personName:       personName,
	}

	event.meta = es.BuildEventMeta(event, streamVersion)

	return event
}

func RebuildCustomerRegistered(
	customerID values.CustomerID,
	emailAddress values.EmailAddress,
	confirmationHash values.ConfirmationHash,
	personName values.PersonName,
	meta es.EventMeta,
) CustomerRegistered {

	event := CustomerRegistered{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
		personName:       personName,
		meta:             meta,
	}

	return event
}

func (event CustomerRegistered) CustomerID() values.CustomerID {
	return event.customerID
}

func (event CustomerRegistered) EmailAddress() values.EmailAddress {
	return event.emailAddress
}

func (event CustomerRegistered) ConfirmationHash() values.ConfirmationHash {
	return event.confirmationHash
}

func (event CustomerRegistered) PersonName() values.PersonName {
	return event.personName
}

func (event CustomerRegistered) Meta() es.EventMeta {
	return event.meta
}

func (event CustomerRegistered) IndicatesAnError() (bool, string) {
	return false, ""
}
