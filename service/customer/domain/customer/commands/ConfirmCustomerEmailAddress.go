package commands

import (
	"go-iddd/service/customer/domain/customer/values"
	"go-iddd/service/lib"
	"reflect"
	"strings"

	"github.com/cockroachdb/errors"
)

type ConfirmCustomerEmailAddress struct {
	customerID       values.CustomerID
	confirmationHash values.ConfirmationHash
	isValid          bool
}

func BuildConfirmCustomerEmailAddress(
	customerID string,
	confirmationHash string,
) (ConfirmCustomerEmailAddress, error) {

	customerIDValue, err := values.BuildCustomerID(customerID)
	if err != nil {
		return ConfirmCustomerEmailAddress{}, err
	}

	confirmationHashValue, err := values.BuildConfirmationHash(confirmationHash)
	if err != nil {
		return ConfirmCustomerEmailAddress{}, err
	}

	confirmEmailAddress := ConfirmCustomerEmailAddress{
		customerID:       customerIDValue,
		confirmationHash: confirmationHashValue,
		isValid:          true,
	}

	return confirmEmailAddress, nil
}

func (command ConfirmCustomerEmailAddress) CustomerID() values.CustomerID {
	return command.customerID
}

func (command ConfirmCustomerEmailAddress) ConfirmationHash() values.ConfirmationHash {
	return command.confirmationHash
}

func (command ConfirmCustomerEmailAddress) ShouldBeValid() error {
	if !command.isValid {
		err := errors.Newf("%s: is not valid", command.commandName())

		return errors.Mark(err, lib.ErrCommandIsInvalid)
	}

	return nil
}

func (command ConfirmCustomerEmailAddress) commandName() string {
	commandType := reflect.TypeOf(command).String()
	commandTypeParts := strings.Split(commandType, ".")
	commandName := commandTypeParts[len(commandTypeParts)-1]

	return strings.Title(commandName)
}
