package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type RegisterCustomer struct {
	customerID       value.CustomerID
	emailAddress     value.EmailAddress
	confirmationHash value.ConfirmationHash
	personName       value.PersonName
	messageID        es.MessageID
}

func BuildRegisterCustomer(
	customerID value.CustomerID,
	emailAddress value.EmailAddress,
	personName value.PersonName,
) RegisterCustomer {

	register := RegisterCustomer{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: value.GenerateConfirmationHash(emailAddress.String()),
		personName:       personName,
		messageID:        es.GenerateMessageID(),
	}

	return register
}

func (command RegisterCustomer) CustomerID() value.CustomerID {
	return command.customerID
}

func (command RegisterCustomer) EmailAddress() value.EmailAddress {
	return command.emailAddress
}

func (command RegisterCustomer) ConfirmationHash() value.ConfirmationHash {
	return command.confirmationHash
}

func (command RegisterCustomer) PersonName() value.PersonName {
	return command.personName
}

func (command RegisterCustomer) MessageID() es.MessageID {
	return command.messageID
}
