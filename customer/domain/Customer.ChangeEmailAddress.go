package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
)

func (customer *customer) ChangeEmailAddress(with *commands.ChangeEmailAddress) {
	if customer.confirmableEmailAddress.Equals(with.EmailAddress()) {
		return
	}

	customer.recordThat(
		events.EmailAddressWasChanged(
			with.ID(),
			with.EmailAddress().ToConfirmable(),
			customer.currentStreamVersion+1,
		),
	)
}

func (customer *customer) whenEmailAddressWasChanged(actualEvent *events.EmailAddressChanged) {
	customer.confirmableEmailAddress = actualEvent.ConfirmableEmailAddress()
}
