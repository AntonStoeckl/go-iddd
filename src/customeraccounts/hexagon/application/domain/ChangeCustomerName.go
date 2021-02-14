package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
)

type ChangeCustomerName struct {
	customerID value.CustomerID
	personName value.PersonName
	messageID  es.MessageID
}

func BuildChangeCustomerName(
	customerID value.CustomerID,
	personName value.PersonName,
) ChangeCustomerName {

	changeEmailAddress := ChangeCustomerName{
		customerID: customerID,
		personName: personName,
		messageID:  es.GenerateMessageID(),
	}

	return changeEmailAddress
}

func (command ChangeCustomerName) CustomerID() value.CustomerID {
	return command.customerID
}

func (command ChangeCustomerName) PersonName() value.PersonName {
	return command.personName
}

func (command ChangeCustomerName) MessageID() es.MessageID {
	return command.messageID
}
