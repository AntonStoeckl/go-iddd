package commands

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
)

type ChangeCustomerEmailAddress struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
}

func BuildChangeCustomerEmailAddress(
	customerID string,
	emailAddress string,
) (ChangeCustomerEmailAddress, error) {

	customerIDValue, err := values.BuildCustomerID(customerID)
	if err != nil {
		return ChangeCustomerEmailAddress{}, err
	}

	emailAddressValue, err := values.BuildEmailAddress(emailAddress)
	if err != nil {
		return ChangeCustomerEmailAddress{}, err
	}

	changeEmailAddress := ChangeCustomerEmailAddress{
		customerID:       customerIDValue,
		emailAddress:     emailAddressValue,
		confirmationHash: values.GenerateConfirmationHash(emailAddressValue.String()),
	}

	return changeEmailAddress, nil
}

func (command ChangeCustomerEmailAddress) CustomerID() values.CustomerID {
	return command.customerID
}

func (command ChangeCustomerEmailAddress) EmailAddress() values.EmailAddress {
	return command.emailAddress
}

func (command ChangeCustomerEmailAddress) ConfirmationHash() values.ConfirmationHash {
	return command.confirmationHash
}
