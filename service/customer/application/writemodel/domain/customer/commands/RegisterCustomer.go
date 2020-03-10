package commands

import (
	"go-iddd/service/customer/application/writemodel/domain/customer/values"
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

func (register RegisterCustomer) CustomerID() values.CustomerID {
	return register.customerID
}

func (register RegisterCustomer) EmailAddress() values.EmailAddress {
	return register.emailAddress
}

func (register RegisterCustomer) ConfirmationHash() values.ConfirmationHash {
	return register.confirmationHash
}

func (register RegisterCustomer) PersonName() values.PersonName {
	return register.personName
}

func (register RegisterCustomer) ShouldBeValid() error {
	if !register.isValid {
		err := errors.Newf("%s: is not valid", register.commandName())

		return errors.Mark(err, lib.ErrCommandIsInvalid)
	}

	return nil
}

func (register RegisterCustomer) commandName() string {
	commandType := reflect.TypeOf(register).String()
	commandTypeParts := strings.Split(commandType, ".")
	commandName := commandTypeParts[len(commandTypeParts)-1]

	return strings.Title(commandName)
}
