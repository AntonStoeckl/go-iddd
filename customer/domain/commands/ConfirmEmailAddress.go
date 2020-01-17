package commands

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"reflect"
	"strings"
)

type ConfirmEmailAddress struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
}

func NewConfirmEmailAddress(
	customerID string,
	emailAddress string,
	confirmationHash string,
) (ConfirmEmailAddress, error) {

	customerIDValue, err := values.BuildCustomerID(customerID)
	if err != nil {
		return ConfirmEmailAddress{}, err
	}

	emailAddressValue, err := values.BuildEmailAddress(emailAddress)
	if err != nil {
		return ConfirmEmailAddress{}, err
	}

	confirmationHashValue, err := values.BuildConfirmationHash(confirmationHash)
	if err != nil {
		return ConfirmEmailAddress{}, err
	}

	confirmEmailAddress := ConfirmEmailAddress{
		customerID:       customerIDValue,
		emailAddress:     emailAddressValue,
		confirmationHash: confirmationHashValue,
	}

	return confirmEmailAddress, nil
}

func (confirmEmailAddress ConfirmEmailAddress) CustomerID() values.CustomerID {
	return confirmEmailAddress.customerID
}

func (confirmEmailAddress ConfirmEmailAddress) EmailAddress() values.EmailAddress {
	return confirmEmailAddress.emailAddress
}

func (confirmEmailAddress ConfirmEmailAddress) ConfirmationHash() values.ConfirmationHash {
	return confirmEmailAddress.confirmationHash
}

func (confirmEmailAddress ConfirmEmailAddress) AggregateID() shared.IdentifiesAggregates {
	return confirmEmailAddress.customerID
}

func (confirmEmailAddress ConfirmEmailAddress) CommandName() string {
	commandType := reflect.TypeOf(confirmEmailAddress).String()
	commandTypeParts := strings.Split(commandType, ".")
	commandName := commandTypeParts[len(commandTypeParts)-1]

	return strings.Title(commandName)
}
