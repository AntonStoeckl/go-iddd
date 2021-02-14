package customer

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
)

func Register(command domain.RegisterCustomer) domain.CustomerRegistered {
	event := domain.BuildCustomerRegistered(
		command.CustomerID(),
		command.EmailAddress(),
		command.ConfirmationHash(),
		command.PersonName(),
		command.MessageID(),
		1,
	)

	return event
}
