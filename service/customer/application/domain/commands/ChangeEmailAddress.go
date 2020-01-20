package commands

import (
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib"
	"reflect"
	"strings"

	"github.com/cockroachdb/errors"
)

type ChangeEmailAddress struct {
	customerID   values.CustomerID
	emailAddress values.EmailAddress
	isValid      bool
}

func NewChangeEmailAddress(
	customerID string,
	emailAddress string,
) (ChangeEmailAddress, error) {

	customerIDValue, err := values.BuildCustomerID(customerID)
	if err != nil {
		return ChangeEmailAddress{}, err
	}

	emailAddressValue, err := values.BuildEmailAddress(emailAddress)
	if err != nil {
		return ChangeEmailAddress{}, err
	}

	changeEmailAddress := ChangeEmailAddress{
		customerID:   customerIDValue,
		emailAddress: emailAddressValue,
		isValid:      true,
	}

	return changeEmailAddress, nil
}

func (changeEmailAddress ChangeEmailAddress) CustomerID() values.CustomerID {
	return changeEmailAddress.customerID
}

func (changeEmailAddress ChangeEmailAddress) EmailAddress() values.EmailAddress {
	return changeEmailAddress.emailAddress
}

func (changeEmailAddress ChangeEmailAddress) AggregateID() lib.IdentifiesAggregates {
	return changeEmailAddress.customerID
}

func (changeEmailAddress ChangeEmailAddress) ShouldBeValid() error {
	if !changeEmailAddress.isValid {
		err := errors.Newf("%s: is not valid", changeEmailAddress.commandName())

		return errors.Mark(err, lib.ErrCommandIsInvalid)
	}

	return nil
}

func (changeEmailAddress ChangeEmailAddress) commandName() string {
	commandType := reflect.TypeOf(changeEmailAddress).String()
	commandTypeParts := strings.Split(commandType, ".")
	commandName := commandTypeParts[len(commandTypeParts)-1]

	return strings.Title(commandName)
}
