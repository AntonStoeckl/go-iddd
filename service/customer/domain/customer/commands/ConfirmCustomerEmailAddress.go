package commands

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
)

type ConfirmCustomerEmailAddress struct {
	customerID       values.CustomerID
	confirmationHash values.ConfirmationHash
}

func BuildConfirmCustomerEmailAddress(
	customerID string,
	confirmationHash string,
) (ConfirmCustomerEmailAddress, error) {

	customerIDValue, err := values.BuildCustomerID(customerID)
	if err != nil {
		return ConfirmCustomerEmailAddress{}, err
	}

	confirmationHashValue, err := values.BuildConfirmationHash(confirmationHash)
	if err != nil {
		return ConfirmCustomerEmailAddress{}, err
	}

	confirmEmailAddress := ConfirmCustomerEmailAddress{
		customerID:       customerIDValue,
		confirmationHash: confirmationHashValue,
	}

	return confirmEmailAddress, nil
}

func (command ConfirmCustomerEmailAddress) CustomerID() values.CustomerID {
	return command.customerID
}

func (command ConfirmCustomerEmailAddress) ConfirmationHash() values.ConfirmationHash {
	return command.confirmationHash
}
