package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/shared"
)

func (customer *Customer) ConfirmEmailAddress(with *commands.ConfirmEmailAddress) shared.DomainEvents {
	if customer.confirmableEmailAddress.IsConfirmed() {
		return nil
	}

	if !customer.confirmableEmailAddress.CanBeConfirmed(with.ConfirmationHash()) {
		event := events.EmailAddressConfirmationHasFailed(
			with.ID(),
			with.ConfirmationHash(),
			customer.currentStreamVersion+1,
		)

		customer.apply(event)

		return shared.DomainEvents{event}
	}

	event := events.EmailAddressWasConfirmed(
		with.ID(),
		with.EmailAddress(),
		customer.currentStreamVersion+1,
	)

	customer.apply(event)

	return shared.DomainEvents{event}
}
