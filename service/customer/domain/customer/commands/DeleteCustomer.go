package commands

import (
	"reflect"
	"strings"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
)

type DeleteCustomer struct {
	customerID values.CustomerID
	isValid    bool
}

func BuildCDeleteCustomer(
	customerID string,
) (DeleteCustomer, error) {

	customerIDValue, err := values.BuildCustomerID(customerID)
	if err != nil {
		return DeleteCustomer{}, err
	}

	deleteCustomer := DeleteCustomer{
		customerID: customerIDValue,
		isValid:    true,
	}

	return deleteCustomer, nil
}

func (command DeleteCustomer) CustomerID() values.CustomerID {
	return command.customerID
}

func (command DeleteCustomer) ShouldBeValid() error {
	if !command.isValid {
		err := errors.Newf("%s: is not valid", command.commandName())

		return errors.Mark(err, lib.ErrCommandIsInvalid)
	}

	return nil
}

func (command DeleteCustomer) commandName() string {
	commandType := reflect.TypeOf(command).String()
	commandTypeParts := strings.Split(commandType, ".")
	commandName := commandTypeParts[len(commandTypeParts)-1]

	return strings.Title(commandName)
}
