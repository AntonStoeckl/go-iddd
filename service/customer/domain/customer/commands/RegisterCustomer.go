package commands

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
)

type RegisterCustomer struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	personName       values.PersonName
}

func BuildRegisterCustomer(
	emailAddress string,
	givenName string,
	familyName string,
) (RegisterCustomer, error) {

	emailAddressValue, err := values.BuildEmailAddress(emailAddress)
	if err != nil {
		return RegisterCustomer{}, err
	}

	personNameValue, err := values.BuildPersonName(givenName, familyName)
	if err != nil {
		return RegisterCustomer{}, err
	}

	register := RegisterCustomer{
		customerID:       values.GenerateCustomerID(),
		emailAddress:     emailAddressValue,
		confirmationHash: values.GenerateConfirmationHash(emailAddressValue.String()),
		personName:       personNameValue,
	}

	return register, nil
}

func (command RegisterCustomer) CustomerID() values.CustomerID {
	return command.customerID
}

func (command RegisterCustomer) EmailAddress() values.EmailAddress {
	return command.emailAddress
}

func (command RegisterCustomer) ConfirmationHash() values.ConfirmationHash {
	return command.confirmationHash
}

func (command RegisterCustomer) PersonName() values.PersonName {
	return command.personName
}
