package domain

import (
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/lib/es"
)

func RegisterCustomer(with commands.Register) es.DomainEvents {
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
