package commands

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"reflect"
	"strings"

	"github.com/cockroachdb/errors"
)

type ConfirmEmailAddress struct {
	customerID       values.CustomerID
	emailAddress     values.EmailAddress
	confirmationHash values.ConfirmationHash
	isValid          bool
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
		isValid:          true,
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

func (confirmEmailAddress ConfirmEmailAddress) ShouldBeValid() error {
	if !confirmEmailAddress.isValid {
		err := errors.Newf("%s: is not valid", confirmEmailAddress.CommandName())

		return errors.Mark(err, shared.ErrCommandIsInvalid)
	}

	return nil
}
