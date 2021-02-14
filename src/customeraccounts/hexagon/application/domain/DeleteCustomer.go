package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type DeleteCustomer struct {
	customerID value.CustomerID
	messageID  es.MessageID
}

func BuildDeleteCustomer(customerID value.CustomerID) DeleteCustomer {
	deleteCustomer := DeleteCustomer{
		customerID: customerID,
		messageID:  es.GenerateMessageID(),
	}

	return deleteCustomer
}

func (command DeleteCustomer) CustomerID() value.CustomerID {
	return command.customerID
}

func (command DeleteCustomer) MessageID() es.MessageID {
	return command.messageID
}
