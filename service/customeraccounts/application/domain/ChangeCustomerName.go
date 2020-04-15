package domain

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer/value"
)

type ChangeCustomerName struct {
	customerID value.CustomerID
	personName value.PersonName
}

func BuildChangeCustomerName(
	customerID string,
	givenName string,
	familyName string,
) (ChangeCustomerName, error) {

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return ChangeCustomerName{}, err
	}

	personNameValue, err := value.BuildPersonName(givenName, familyName)
	if err != nil {
		return ChangeCustomerName{}, err
	}

	changeEmailAddress := ChangeCustomerName{
		customerID: customerIDValue,
		personName: personNameValue,
	}

	return changeEmailAddress, nil
}

func (command ChangeCustomerName) CustomerID() value.CustomerID {
	return command.customerID
}

func (command ChangeCustomerName) PersonName() value.PersonName {
	return command.personName
}
