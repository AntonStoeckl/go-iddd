package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
)

func NewCustomerWith(register *commands.Register) Customer {
	newCustomer := blankCustomer()

	newCustomer.recordThat(
		events.ItWasRegistered(
			register.ID(),
			register.EmailAddress().ToConfirmable(),
			register.PersonName(),
		),
	)

	return newCustomer
}

func (customer *customer) whenItWasRegistered(actualEvent *events.Registered) {
	customer.id = actualEvent.ID()
	customer.confirmableEmailAddress = actualEvent.ConfirmableEmailAddress()
	customer.personName = actualEvent.PersonName()
}
