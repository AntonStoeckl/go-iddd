package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
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
	emailAddress string,
	givenName string,
	familyName string,
) (RegisterCustomer, error) {

	wrapWithMsg := "BuildRegisterCustomer"

	emailAddressValue, err := value.BuildEmailAddress(emailAddress)
	if err != nil {
		return RegisterCustomer{}, errors.Wrap(err, wrapWithMsg)
	}

	personNameValue, err := value.BuildPersonName(givenName, familyName)
	if err != nil {
		return RegisterCustomer{}, errors.Wrap(err, wrapWithMsg)
	}

	command := RegisterCustomer{
		customerID:       customerID,
		emailAddress:     emailAddressValue,
		confirmationHash: value.GenerateConfirmationHash(emailAddressValue.String()),
		personName:       personNameValue,
		messageID:        es.GenerateMessageID(),
	}

	return command, nil
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
