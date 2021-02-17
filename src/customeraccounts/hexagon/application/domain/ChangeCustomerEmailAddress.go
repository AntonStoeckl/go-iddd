package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
)

type ChangeCustomerEmailAddress struct {
	customerID       value.CustomerID
	emailAddress     value.UnconfirmedEmailAddress
	confirmationHash value.ConfirmationHash
	messageID        es.MessageID
}

func BuildChangeCustomerEmailAddress(
	customerID string,
	emailAddress string,
) (ChangeCustomerEmailAddress, error) {

	wrapWithMsg := "customerCommandHandler.ChangeCustomerEmailAddress"

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return ChangeCustomerEmailAddress{}, errors.Wrap(err, wrapWithMsg)
	}

	emailAddressValue, err := value.BuildUnconfirmedEmailAddress(emailAddress)
	if err != nil {
		return ChangeCustomerEmailAddress{}, errors.Wrap(err, wrapWithMsg)
	}

	changeEmailAddress := ChangeCustomerEmailAddress{
		customerID:       customerIDValue,
		emailAddress:     emailAddressValue,
		confirmationHash: value.GenerateConfirmationHash(emailAddressValue.String()),
		messageID:        es.GenerateMessageID(),
	}

	return changeEmailAddress, nil
}

func (command ChangeCustomerEmailAddress) CustomerID() value.CustomerID {
	return command.customerID
}

func (command ChangeCustomerEmailAddress) EmailAddress() value.UnconfirmedEmailAddress {
	return command.emailAddress
}

func (command ChangeCustomerEmailAddress) ConfirmationHash() value.ConfirmationHash {
	return command.confirmationHash
}

func (command ChangeCustomerEmailAddress) MessageID() es.MessageID {
	return command.messageID
}
