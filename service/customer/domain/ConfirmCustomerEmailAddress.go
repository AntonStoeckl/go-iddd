package domain

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/value"
)

type ConfirmCustomerEmailAddress struct {
	customerID       value.CustomerID
	confirmationHash value.ConfirmationHash
}

func BuildConfirmCustomerEmailAddress(
	customerID string,
	confirmationHash string,
) (ConfirmCustomerEmailAddress, error) {

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return ConfirmCustomerEmailAddress{}, err
	}

	confirmationHashValue, err := value.BuildConfirmationHash(confirmationHash)
	if err != nil {
		return ConfirmCustomerEmailAddress{}, err
	}

	confirmEmailAddress := ConfirmCustomerEmailAddress{
		customerID:       customerIDValue,
		confirmationHash: confirmationHashValue,
	}

	return confirmEmailAddress, nil
}

func (command ConfirmCustomerEmailAddress) CustomerID() value.CustomerID {
	return command.customerID
}

func (command ConfirmCustomerEmailAddress) ConfirmationHash() value.ConfirmationHash {
	return command.confirmationHash
}
