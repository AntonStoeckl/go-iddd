package commands

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
)

type ChangeCustomerName struct {
	customerID values.CustomerID
	personName values.PersonName
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
	}

	return changeEmailAddress, nil
}

func (command ChangeCustomerName) CustomerID() values.CustomerID {
	return command.customerID
}

func (command ChangeCustomerName) PersonName() values.PersonName {
	return command.personName
}
