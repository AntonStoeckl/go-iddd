package domain

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer/value"
)

type ConfirmCustomerEmailAddress struct {
	customerID       value.CustomerID
	confirmationHash value.ConfirmationHash
}

func BuildConfirmCustomerEmailAddress(
	customerID value.CustomerID,
	confirmationHash value.ConfirmationHash,
) ConfirmCustomerEmailAddress {

	confirmEmailAddress := ConfirmCustomerEmailAddress{
		customerID:       customerID,
		confirmationHash: confirmationHash,
	}

	return confirmEmailAddress
}

func (command ConfirmCustomerEmailAddress) CustomerID() value.CustomerID {
	return command.customerID
}

func (command ConfirmCustomerEmailAddress) ConfirmationHash() value.ConfirmationHash {
	return command.confirmationHash
}
