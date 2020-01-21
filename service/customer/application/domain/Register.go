package domain

import (
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/lib"
)

func RegisterCustomer(with commands.Register) lib.DomainEvents {
	return lib.DomainEvents{
		events.ItWasRegistered(
			with.CustomerID(),
			with.EmailAddress(),
			with.ConfirmationHash(),
			with.PersonName(),
			1,
		),
	}
}
