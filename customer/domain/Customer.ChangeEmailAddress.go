package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/shared"
)

func (customer *Customer) ChangeEmailAddress(with *commands.ChangeEmailAddress) shared.DomainEvents {
	if customer.confirmableEmailAddress.Equals(with.EmailAddress()) {
		return nil
	}

	event := events.EmailAddressWasChanged(
		with.CustomerID(),
		with.EmailAddress().ToConfirmable(),
		customer.currentStreamVersion+1,
	)

	customer.apply(event)

	return shared.DomainEvents{event}
}
