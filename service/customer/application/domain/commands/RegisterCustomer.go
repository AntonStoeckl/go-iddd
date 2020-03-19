package commands

import (
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib"
	"reflect"
	"strings"

	"github.com/cockroachdb/errors"
)

type RegisterCustomer struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	personName       values.PersonName
	isValid          bool
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
		confirmationHash: values.GenerateConfirmationHash(emailAddressValue.EmailAddress()),
		personName:       personNameValue,
		isValid:          true,
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

func (command RegisterCustomer) ShouldBeValid() error {
	if !command.isValid {
		err := errors.Newf("%s: is not valid", command.commandName())

		return errors.Mark(err, lib.ErrCommandIsInvalid)
	}

	return nil
}

func (command RegisterCustomer) commandName() string {
	commandType := reflect.TypeOf(command).String()
	commandTypeParts := strings.Split(commandType, ".")
	commandName := commandTypeParts[len(commandTypeParts)-1]

	return strings.Title(commandName)
}
