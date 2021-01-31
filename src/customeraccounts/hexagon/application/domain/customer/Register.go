package customer

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
)

func Register(with domain.RegisterCustomer) domain.CustomerRegistered {
	event := domain.BuildCustomerRegistered(
		with.CustomerID(),
		with.EmailAddress(),
		with.ConfirmationHash(),
		with.PersonName(),
		1,
	)

	return event
}