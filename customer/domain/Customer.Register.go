package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

func RegisterCustomer(with commands.Register) shared.DomainEvents {
	return shared.DomainEvents{
		events.ItWasRegistered(
			with.CustomerID(),
			with.EmailAddress(),
			values.GenerateConfirmationHash(with.EmailAddress().EmailAddress()),
			with.PersonName(),
			1,
		),
	}
}
