package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
)

type ConfirmCustomerEmailAddress struct {
	customerID       value.CustomerID
	confirmationHash value.ConfirmationHash
	messageID        es.MessageID
}

func BuildConfirmCustomerEmailAddress(
	customerID string,
	confirmationHash string,
) (ConfirmCustomerEmailAddress, error) {

	wrapWithMsg := "BuildConfirmCustomerEmailAddress"

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return ConfirmCustomerEmailAddress{}, errors.Wrap(err, wrapWithMsg)
	}

	confirmationHashValue, err := value.BuildConfirmationHash(confirmationHash)
	if err != nil {
		return ConfirmCustomerEmailAddress{}, errors.Wrap(err, wrapWithMsg)
	}

	command := ConfirmCustomerEmailAddress{
		customerID:       customerIDValue,
		confirmationHash: confirmationHashValue,
		messageID:        es.GenerateMessageID(),
	}

	return command, nil
}

func (command ConfirmCustomerEmailAddress) CustomerID() value.CustomerID {
	return command.customerID
}

func (command ConfirmCustomerEmailAddress) ConfirmationHash() value.ConfirmationHash {
	return command.confirmationHash
}

func (command ConfirmCustomerEmailAddress) MessageID() es.MessageID {
	return command.messageID
}
