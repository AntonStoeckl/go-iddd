package domain

import (
	"go-iddd/service/customer/domain/commands"
	"go-iddd/service/customer/domain/events"
	"go-iddd/service/customer/domain/values"
	"go-iddd/service/lib"
)

func RegisterCustomer(with commands.Register) lib.DomainEvents {
	return lib.DomainEvents{
		events.ItWasRegistered(
			with.CustomerID(),
			with.EmailAddress(),
			values.GenerateConfirmationHash(with.EmailAddress().EmailAddress()),
			with.PersonName(),
			1,
		),
	}
}
