package domain

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer/value"
)

type RegisterCustomer struct {
	customerID       value.CustomerID
	emailAddress     value.EmailAddress
	confirmationHash value.ConfirmationHash
	personName       value.PersonName
}

func BuildRegisterCustomer(
	customerID value.CustomerID,
	emailAddress value.EmailAddress,
	confirmationHash value.ConfirmationHash,
	personName value.PersonName,
) RegisterCustomer {

	register := RegisterCustomer{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: confirmationHash,
		personName:       personName,
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
