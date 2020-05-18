package domain

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer/value"
)

type ChangeCustomerName struct {
	customerID value.CustomerID
	personName value.PersonName
}

func BuildChangeCustomerName(
	customerID value.CustomerID,
	personName value.PersonName,
) ChangeCustomerName {

	changeEmailAddress := ChangeCustomerName{
		customerID: customerID,
		personName: personName,
	}

	return changeEmailAddress
}

func (command ChangeCustomerName) CustomerID() value.CustomerID {
	return command.customerID
}

func (command ChangeCustomerName) PersonName() value.PersonName {
	return command.personName
}
