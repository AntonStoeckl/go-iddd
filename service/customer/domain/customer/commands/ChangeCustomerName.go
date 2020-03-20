package commands

import (
	"go-iddd/service/customer/domain/customer/values"
	"go-iddd/service/lib"
	"reflect"
	"strings"

	"github.com/cockroachdb/errors"
)

type ChangeCustomerName struct {
	customerID values.CustomerID
	personName values.PersonName
	isValid    bool
}

func BuildChangeCustomerName(
	customerID string,
	givenName string,
	familyName string,
) (ChangeCustomerName, error) {

	customerIDValue, err := values.BuildCustomerID(customerID)
	if err != nil {
		return ChangeCustomerName{}, err
	}

	personNameValue, err := values.BuildPersonName(givenName, familyName)
	if err != nil {
		return ChangeCustomerName{}, err
	}

	changeEmailAddress := ChangeCustomerName{
		customerID: customerIDValue,
		personName: personNameValue,
		isValid:    true,
	}

	return changeEmailAddress, nil
}

func (command ChangeCustomerName) CustomerID() values.CustomerID {
	return command.customerID
}

func (command ChangeCustomerName) PersonName() values.PersonName {
	return command.personName
}

func (command ChangeCustomerName) ShouldBeValid() error {
	if !command.isValid {
		err := errors.Newf("%s: is not valid", command.commandName())

		return errors.Mark(err, lib.ErrCommandIsInvalid)
	}

	return nil
}

func (command ChangeCustomerName) commandName() string {
	commandType := reflect.TypeOf(command).String()
	commandTypeParts := strings.Split(commandType, ".")
	commandName := commandTypeParts[len(commandTypeParts)-1]

	return strings.Title(commandName)
}
