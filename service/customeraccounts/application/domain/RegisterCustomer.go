package domain

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer/value"
)

type RegisterCustomer struct {
	customerID       value.CustomerID
	emailAddress     value.EmailAddress
	confirmationHash value.ConfirmationHash
	personName       value.PersonName
}

func BuildRegisterCustomer(
	emailAddress string,
	givenName string,
	familyName string,
) (RegisterCustomer, error) {

	emailAddressValue, err := value.BuildEmailAddress(emailAddress)
	if err != nil {
		return RegisterCustomer{}, err
	}

	personNameValue, err := value.BuildPersonName(givenName, familyName)
	if err != nil {
		return RegisterCustomer{}, err
	}

	register := RegisterCustomer{
		customerID:       value.GenerateCustomerID(),
		emailAddress:     emailAddressValue,
		confirmationHash: value.GenerateConfirmationHash(emailAddressValue.String()),
		personName:       personNameValue,
	}

	return register, nil
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
