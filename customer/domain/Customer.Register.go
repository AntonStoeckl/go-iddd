package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
)

func Register(with *commands.Register) Customer {
	newCustomer := blankCustomer()

	newCustomer.recordThat(
		events.ItWasRegistered(
			with.ID(),
			with.EmailAddress().ToConfirmable(),
			with.PersonName(),
		),
	)

	return newCustomer
}

func (customer *customer) whenItWasRegistered(actualEvent *events.Registered) {
	customer.id = actualEvent.ID()
	customer.confirmableEmailAddress = actualEvent.ConfirmableEmailAddress()
	customer.personName = actualEvent.PersonName()
}
