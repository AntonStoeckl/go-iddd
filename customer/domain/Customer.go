package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"

	"github.com/cockroachdb/errors"
)

type Customer interface {
	ID() shared.IdentifiesAggregates
	ConfirmEmailAddress(with *commands.ConfirmEmailAddress) error
	ChangeEmailAddress(with *commands.ChangeEmailAddress)
	StreamVersion() uint
	RecordedEvents(purge bool) shared.DomainEvents
}

type customer struct {
	id                      *values.ID
	confirmableEmailAddress *values.ConfirmableEmailAddress
	personName              *values.PersonName
	currentStreamVersion    uint
	recordedEvents          shared.DomainEvents
}

func blankCustomer() *customer {
	return &customer{}
}

func (customer *customer) ID() shared.IdentifiesAggregates {
	return customer.id
}

func ReconstituteCustomerFrom(eventStream shared.DomainEvents) (Customer, error) {
	newCustomer := blankCustomer()

	if err := eventStream.FirstEventShouldBeOfSameTypeAs(&events.Registered{}); err != nil {
		return nil, errors.Wrap(errors.Mark(err, shared.ErrInvalidEventStream), "shared.ErrInvalidEventStream")
	}

	newCustomer.apply(eventStream...)

	return newCustomer, nil
}

func (customer *customer) recordThat(event shared.DomainEvent) {
	customer.recordedEvents = append(customer.recordedEvents, event)
	customer.apply(event)
}

func (customer *customer) apply(eventStream ...shared.DomainEvent) {
	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case *events.Registered:
			customer.whenItWasRegistered(actualEvent)
		case *events.EmailAddressConfirmed:
			customer.whenEmailAddressWasConfirmed(actualEvent)
		case *events.EmailAddressChanged:
			customer.whenEmailAddressWasChanged(actualEvent)
		}

		customer.currentStreamVersion = event.StreamVersion()
	}
}

func (customer *customer) RecordedEvents(purge bool) shared.DomainEvents {
	recordedEvents := customer.recordedEvents

	if purge {
		customer.recordedEvents = nil
	}

	return recordedEvents
}

func (customer *customer) StreamVersion() uint {
	return customer.currentStreamVersion
}
