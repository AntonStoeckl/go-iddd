package commands

import (
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib"
	"reflect"
	"strings"

	"github.com/cockroachdb/errors"
)

type Register struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	personName       values.PersonName
	isValid          bool
}

func NewRegister(
	emailAddress string,
	givenName string,
	familyName string,
) (Register, error) {

	emailAddressValue, err := values.BuildEmailAddress(emailAddress)
	if err != nil {
		return Register{}, err
	}

	confirmationHash := values.GenerateConfirmationHash(emailAddressValue.EmailAddress())

	personNameValue, err := values.BuildPersonName(givenName, familyName)
	if err != nil {
		return Register{}, err
	}

	register := Register{
		customerID:       values.GenerateCustomerID(),
		emailAddress:     emailAddressValue,
		confirmationHash: confirmationHash,
		personName:       personNameValue,
		isValid:          true,
	}

	return register, nil
}

func (register Register) CustomerID() values.CustomerID {
	return register.customerID
}

func (register Register) EmailAddress() values.EmailAddress {
	return register.emailAddress
}

func (register Register) ConfirmationHash() values.ConfirmationHash {
	return register.confirmationHash
}

func (register Register) PersonName() values.PersonName {
	return register.personName
}

func (register Register) ShouldBeValid() error {
	if !register.isValid {
		err := errors.Newf("%s: is not valid", register.commandName())

		return errors.Mark(err, lib.ErrCommandIsInvalid)
	}

	return nil
}

func (register Register) commandName() string {
	commandType := reflect.TypeOf(register).String()
	commandTypeParts := strings.Split(commandType, ".")
	commandName := commandTypeParts[len(commandTypeParts)-1]

	return strings.Title(commandName)
}

func (register Register) IsCommand() bool {
	return true
}
