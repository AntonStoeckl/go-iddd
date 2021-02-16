package domain

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
)

type ChangeCustomerName struct {
	customerID value.CustomerID
	personName value.PersonName
	messageID  es.MessageID
}

func BuildChangeCustomerName(
	customerID string,
	givenName string,
	familyName string,
) (ChangeCustomerName, error) {

	wrapWithMsg := "BuildChangeCustomerName"

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return ChangeCustomerName{}, errors.Wrap(err, wrapWithMsg)
	}

	personNameValue, err := value.BuildPersonName(givenName, familyName)
	if err != nil {
		return ChangeCustomerName{}, errors.Wrap(err, wrapWithMsg)
	}

	command := ChangeCustomerName{
		customerID: customerIDValue,
		personName: personNameValue,
		messageID:  es.GenerateMessageID(),
	}

	return command, nil
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
