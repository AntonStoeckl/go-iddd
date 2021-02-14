package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type ChangeCustomerEmailAddress struct {
	customerID       value.CustomerID
	emailAddress     value.EmailAddress
	confirmationHash value.ConfirmationHash
	messageID        es.MessageID
}

func BuildChangeCustomerEmailAddress(
	customerID value.CustomerID,
	emailAddress value.EmailAddress,
) ChangeCustomerEmailAddress {

	changeEmailAddress := ChangeCustomerEmailAddress{
		customerID:       customerID,
		emailAddress:     emailAddress,
		confirmationHash: value.GenerateConfirmationHash(emailAddress.String()),
		messageID:        es.GenerateMessageID(),
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

func (command ChangeCustomerEmailAddress) MessageID() es.MessageID {
	return command.messageID
}
