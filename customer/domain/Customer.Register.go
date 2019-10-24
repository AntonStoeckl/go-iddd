package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/shared"
)

func RegisterCustomer(with *commands.Register) shared.DomainEvents {
	return shared.DomainEvents{
		events.ItWasRegistered(
			with.ID(),
			with.EmailAddress().ToConfirmable(),
			with.PersonName(),
			1,
		),
	}
}
