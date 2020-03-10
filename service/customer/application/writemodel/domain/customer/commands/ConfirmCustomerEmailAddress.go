package commands

import (
	"go-iddd/service/customer/application/writemodel/domain/customer/values"
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

func (confirmEmailAddress ConfirmCustomerEmailAddress) CustomerID() values.CustomerID {
	return confirmEmailAddress.customerID
}

func (confirmEmailAddress ConfirmCustomerEmailAddress) ConfirmationHash() values.ConfirmationHash {
	return confirmEmailAddress.confirmationHash
}

func (confirmEmailAddress ConfirmCustomerEmailAddress) ShouldBeValid() error {
	if !confirmEmailAddress.isValid {
		err := errors.Newf("%s: is not valid", confirmEmailAddress.commandName())

		return errors.Mark(err, lib.ErrCommandIsInvalid)
	}

	return nil
}

func (confirmEmailAddress ConfirmCustomerEmailAddress) commandName() string {
	commandType := reflect.TypeOf(confirmEmailAddress).String()
	commandTypeParts := strings.Split(commandType, ".")
	commandName := commandTypeParts[len(commandTypeParts)-1]

	return strings.Title(commandName)
}
