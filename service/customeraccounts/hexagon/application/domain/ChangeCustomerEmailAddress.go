package domain

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer/value"
)

type ChangeCustomerEmailAddress struct {
	customerID       value.CustomerID
	emailAddress     value.EmailAddress
	confirmationHash value.ConfirmationHash
}

func BuildChangeCustomerEmailAddress(
	customerID value.CustomerID,
	emailAddress value.EmailAddress,
) ChangeCustomerEmailAddress {

	changeEmailAddress := ChangeCustomerEmailAddress{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: value.GenerateConfirmationHash(emailAddress.String()),
	}

	return changeEmailAddress
}

func (command ChangeCustomerEmailAddress) CustomerID() value.CustomerID {
	return command.customerID
}

func (command ChangeCustomerEmailAddress) EmailAddress() value.EmailAddress {
	return command.emailAddress
}

func (command ChangeCustomerEmailAddress) ConfirmationHash() value.ConfirmationHash {
	return command.confirmationHash
}
