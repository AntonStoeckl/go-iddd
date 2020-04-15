package domain

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer/value"
)

type DeleteCustomer struct {
	customerID value.CustomerID
}

func BuildDeleteCustomer(
	customerID string,
) (DeleteCustomer, error) {

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return DeleteCustomer{}, err
	}

	deleteCustomer := DeleteCustomer{
		customerID: customerIDValue,
	}

	return deleteCustomer, nil
}

func (command DeleteCustomer) CustomerID() value.CustomerID {
	return command.customerID
}
