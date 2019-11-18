package domain

import (
	"fmt"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"reflect"

	"github.com/cockroachdb/errors"
)

type Customer struct {
	id                      *values.CustomerID
	confirmableEmailAddress *values.ConfirmableEmailAddress
	personName              *values.PersonName
	currentStreamVersion    uint
}

func blankCustomer() *Customer {
	return &Customer{}
}

func (customer *Customer) ID() shared.IdentifiesAggregates {
	return customer.id
}

func ReconstituteCustomerFrom(eventStream shared.DomainEvents) (*Customer, error) {
	newCustomer := blankCustomer()

	if err := newCustomer.shouldBeRegistered(eventStream); err != nil {
		return nil, errors.Wrap(errors.Mark(err, shared.ErrInvalidEventStream), "shared.ErrInvalidEventStream")
	}

	newCustomer.apply(eventStream...)

	return newCustomer, nil
}

func (customer *Customer) shouldBeRegistered(eventStream shared.DomainEvents) error {
	if len(eventStream) == 0 {
		return errors.New("eventStream is empty")
	}

	expectedType := reflect.TypeOf(&events.Registered{})
	actualType := reflect.TypeOf(eventStream[0])

	if actualType != expectedType {
		return fmt.Errorf(
			"first event in eventStream should be [%s] but is type [%s]",
			expectedType.String(),
			actualType.String(),
		)
	}

	return nil
}

func (customer *Customer) apply(eventStream ...shared.DomainEvent) {
	for _, event := range eventStream {
		switch actualEvent := event.(type) {
		case *events.Registered:
			customer.id = actualEvent.ID()
			customer.confirmableEmailAddress = actualEvent.ConfirmableEmailAddress()
			customer.personName = actualEvent.PersonName()
		case *events.EmailAddressConfirmed:
			customer.confirmableEmailAddress = customer.confirmableEmailAddress.MarkAsConfirmed()
		case *events.EmailAddressChanged:
			customer.confirmableEmailAddress = actualEvent.ConfirmableEmailAddress()
		}

		customer.currentStreamVersion = event.StreamVersion()
	}
}

func (customer *Customer) StreamVersion() uint {
	return customer.currentStreamVersion
}
