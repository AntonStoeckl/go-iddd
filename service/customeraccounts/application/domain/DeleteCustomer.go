package domain

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer/value"
)

type DeleteCustomer struct {
	customerID value.CustomerID
}

func BuildDeleteCustomer(
	customerID value.CustomerID,
) DeleteCustomer {

	deleteCustomer := DeleteCustomer{
		customerID: customerID,
	}

	return deleteCustomer
}

func (command DeleteCustomer) CustomerID() value.CustomerID {
	return command.customerID
}
