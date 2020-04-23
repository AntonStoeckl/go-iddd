package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain"
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
