package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type ConfirmCustomerEmailAddress struct {
	customerID       value.CustomerID
	confirmationHash value.ConfirmationHash
	messageID        es.MessageID
}

func BuildConfirmCustomerEmailAddress(
	customerID value.CustomerID,
	confirmationHash value.ConfirmationHash,
) ConfirmCustomerEmailAddress {

	command := ConfirmCustomerEmailAddress{
		customerID:       customerID,
		confirmationHash: confirmationHash,
		messageID:        es.GenerateMessageID(),
	}

	return command
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
