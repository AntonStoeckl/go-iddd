package customer

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
)

func Register(with domain.RegisterCustomer) es.RecordedEvents {
	return es.RecordedEvents{
		domain.BuildCustomerRegistered(
			with.CustomerID(),
			with.EmailAddress(),
			with.ConfirmationHash(),
			with.PersonName(),
			1,
		),
	}
}
