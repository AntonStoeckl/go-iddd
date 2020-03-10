package commands

import (
	"go-iddd/service/customer/application/writemodel/domain/customer/values"
	"go-iddd/service/lib"
	"reflect"
	"strings"

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

func (changeEmailAddress ChangeCustomerEmailAddress) CustomerID() values.CustomerID {
	return changeEmailAddress.customerID
}

func (changeEmailAddress ChangeCustomerEmailAddress) EmailAddress() values.EmailAddress {
	return changeEmailAddress.emailAddress
}

func (changeEmailAddress ChangeCustomerEmailAddress) ConfirmationHash() values.ConfirmationHash {
	return changeEmailAddress.confirmationHash
}

func (changeEmailAddress ChangeCustomerEmailAddress) ShouldBeValid() error {
	if !changeEmailAddress.isValid {
		err := errors.Newf("%s: is not valid", changeEmailAddress.commandName())

		return errors.Mark(err, lib.ErrCommandIsInvalid)
	}

	return nil
}

func (changeEmailAddress ChangeCustomerEmailAddress) commandName() string {
	commandType := reflect.TypeOf(changeEmailAddress).String()
	commandTypeParts := strings.Split(commandType, ".")
	commandName := commandTypeParts[len(commandTypeParts)-1]

	return strings.Title(commandName)
}
