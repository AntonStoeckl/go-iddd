package commands

import (
	"reflect"
	"strings"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
)

type ChangeCustomerEmailAddress struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	isValid          bool
}

func BuildChangeCustomerEmailAddress(
	customerID string,
	emailAddress string,
) (ChangeCustomerEmailAddress, error) {

	customerIDValue, err := values.BuildCustomerID(customerID)
	if err != nil {
		return ChangeCustomerEmailAddress{}, err
	}

	emailAddressValue, err := values.BuildEmailAddress(emailAddress)
	if err != nil {
		return ChangeCustomerEmailAddress{}, err
	}

	changeEmailAddress := ChangeCustomerEmailAddress{
		customerID:       customerIDValue,
		emailAddress:     emailAddressValue,
		confirmationHash: values.GenerateConfirmationHash(emailAddressValue.EmailAddress()),
		isValid:          true,
	}

	return changeEmailAddress, nil
}

func (command ChangeCustomerEmailAddress) CustomerID() values.CustomerID {
	return command.customerID
}

func (command ChangeCustomerEmailAddress) EmailAddress() values.EmailAddress {
	return command.emailAddress
}

func (command ChangeCustomerEmailAddress) ConfirmationHash() values.ConfirmationHash {
	return command.confirmationHash
}

func (command ChangeCustomerEmailAddress) ShouldBeValid() error {
	if !command.isValid {
		err := errors.Newf("%s: is not valid", command.commandName())

		return errors.Mark(err, lib.ErrCommandIsInvalid)
	}

	return nil
}

func (command ChangeCustomerEmailAddress) commandName() string {
	commandType := reflect.TypeOf(command).String()
	commandTypeParts := strings.Split(commandType, ".")
	commandName := commandTypeParts[len(commandTypeParts)-1]

	return strings.Title(commandName)
}
