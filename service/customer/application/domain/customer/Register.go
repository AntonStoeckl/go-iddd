package customer

import (
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/lib/es"
)

func Register(with commands.RegisterCustomer) es.DomainEvents {
	return es.DomainEvents{
		events.CustomerWasRegistered(
			with.CustomerID(),
			with.EmailAddress(),
			with.ConfirmationHash(),
			with.PersonName(),
			1,
		),
	}
}