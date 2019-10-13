package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
)

func RegisterCustomer(with *commands.Register) Customer {
	newCustomer := blankCustomer()

	newCustomer.recordThat(
		events.ItWasRegistered(
			with.ID(),
			with.EmailAddress().ToConfirmable(),
			with.PersonName(),
			newCustomer.currentStreamVersion+1,
		),
	)

	return newCustomer
}

func (customer *customer) whenItWasRegistered(actualEvent *events.Registered) {
	customer.id = actualEvent.ID()
	customer.confirmableEmailAddress = actualEvent.ConfirmableEmailAddress()
	customer.personName = actualEvent.PersonName()
}
