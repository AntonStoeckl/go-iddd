package domain

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer/value"
)

type ChangeCustomerEmailAddress struct {
	customerID       value.CustomerID
	emailAddress     value.EmailAddress
	confirmationHash value.ConfirmationHash
}

func BuildChangeCustomerEmailAddress(
	customerID string,
	emailAddress string,
) (ChangeCustomerEmailAddress, error) {

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return ChangeCustomerEmailAddress{}, err
	}

	emailAddressValue, err := value.BuildEmailAddress(emailAddress)
	if err != nil {
		return ChangeCustomerEmailAddress{}, err
	}

	changeEmailAddress := ChangeCustomerEmailAddress{
		customerID:       customerIDValue,
		emailAddress:     emailAddressValue,
		confirmationHash: value.GenerateConfirmationHash(emailAddressValue.String()),
	}

	return changeEmailAddress, nil
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
