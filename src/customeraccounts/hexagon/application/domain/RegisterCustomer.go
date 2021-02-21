package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type RegisterCustomer struct {
	customerID   value.CustomerID
	emailAddress value.UnconfirmedEmailAddress
	personName   value.PersonName
	messageID    es.MessageID
}

func BuildRegisterCustomer(
	customerID value.CustomerID,
	emailAddress value.UnconfirmedEmailAddress,
	personName value.PersonName,
) RegisterCustomer {

	command := RegisterCustomer{
		customerID:   customerID,
		emailAddress: emailAddress,
		personName:   personName,
		messageID:    es.GenerateMessageID(),
	}

	return command
}

func (command RegisterCustomer) CustomerID() value.CustomerID {
	return command.customerID
}

func (command RegisterCustomer) EmailAddress() value.UnconfirmedEmailAddress {
	return command.emailAddress
}

func (command RegisterCustomer) PersonName() value.PersonName {
	return command.personName
}

func (command RegisterCustomer) MessageID() es.MessageID {
	return command.messageID
}
