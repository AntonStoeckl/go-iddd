package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type ChangeCustomerEmailAddress struct {
	customerID   value.CustomerID
	emailAddress value.UnconfirmedEmailAddress
	messageID    es.MessageID
}

func BuildChangeCustomerEmailAddress(
	customerID value.CustomerID,
	emailAddress value.UnconfirmedEmailAddress,
) ChangeCustomerEmailAddress {

	changeEmailAddress := ChangeCustomerEmailAddress{
		customerID:   customerID,
		emailAddress: emailAddress,
		messageID:    es.GenerateMessageID(),
	}

	return changeEmailAddress
}

func (command ChangeCustomerEmailAddress) CustomerID() value.CustomerID {
	return command.customerID
}

func (command ChangeCustomerEmailAddress) EmailAddress() value.UnconfirmedEmailAddress {
	return command.emailAddress
}

func (command ChangeCustomerEmailAddress) MessageID() es.MessageID {
	return command.messageID
}
