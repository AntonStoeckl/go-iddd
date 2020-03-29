package commands

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
)

type DeleteCustomer struct {
	customerID values.CustomerID
}

func BuildDeleteCustomer(
	customerID string,
) (DeleteCustomer, error) {

	customerIDValue, err := values.BuildCustomerID(customerID)
	if err != nil {
		return DeleteCustomer{}, err
	}

	deleteCustomer := DeleteCustomer{
		customerID: customerIDValue,
	}

	return deleteCustomer, nil
}

func (command DeleteCustomer) CustomerID() values.CustomerID {
	return command.customerID
}
