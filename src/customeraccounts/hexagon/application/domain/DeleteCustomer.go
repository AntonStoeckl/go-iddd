package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
)

type DeleteCustomer struct {
	customerID value.CustomerID
	messageID  es.MessageID
}

func BuildDeleteCustomer(customerID string) (DeleteCustomer, error) {
	wrapWithMsg := "customerCommandHandler.DeleteCustomer"

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return DeleteCustomer{}, errors.Wrap(err, wrapWithMsg)
	}

	command := DeleteCustomer{
		customerID: customerIDValue,
		messageID:  es.GenerateMessageID(),
	}

	return command, nil
}

func (command DeleteCustomer) CustomerID() value.CustomerID {
	return command.customerID
}

func (command DeleteCustomer) MessageID() es.MessageID {
	return command.messageID
}
