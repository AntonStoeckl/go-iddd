package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"reflect"
	"strings"

	"github.com/cockroachdb/errors"
)

type Customer interface {
	Clone() Customer

	shared.EventsourcedAggregate
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

/*** Implement Customer ***/

func (customer *customer) Clone() Customer {
	cloned := *customer

	return &cloned
}

/*** Implement shared.Aggregate ****/

func (customer *customer) Execute(command shared.Command) error {
	var err error

	if err := customer.assertIsValid(command); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrCommandIsInvalid), "customer.Execute")
	}

	switch actualCommand := command.(type) {
	case *commands.ConfirmEmailAddress:
		err = customer.confirmEmailAddress(actualCommand)
	case *commands.ChangeEmailAddress:
		err = customer.changeEmailAddress(actualCommand)
	case *commands.Register:
		err = errors.Mark(errors.New("customer.Execute: customer is already registered"), shared.ErrCommandCanNotBeHandled)
	default:
		err = errors.Mark(errors.Newf("customer.Execute: [%s]: command is unknown", command.CommandName()), shared.ErrCommandCanNotBeHandled)
	}

	if err != nil {
		return err
	}

	return nil
}

func (customer *customer) AggregateID() shared.IdentifiesAggregates {
	return customer.id
}

func (customer *customer) AggregateName() string {
	aggregateType := reflect.TypeOf(customer).String()
	aggregateTypeParts := strings.Split(aggregateType, ".")
	aggregateName := aggregateTypeParts[len(aggregateTypeParts)-1]

	return strings.Title(aggregateName)
}

/*** EventSourcing ***/

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

/*** Implement shared.RecordsEvents ****/

func (customer *customer) RecordedEvents(purge bool) shared.DomainEvents {
	recordedEvents := customer.recordedEvents

	if purge {
		customer.recordedEvents = nil
	}

	return recordedEvents
}

/*** Implement shared.EventsourcedAggregate ****/

func (customer *customer) StreamVersion() uint {
	return customer.currentStreamVersion
}

func (customer *customer) Apply(latestEvents shared.DomainEvents) {
	customer.apply(latestEvents...)
}

/*** Command Assertions ***/

func (customer *customer) assertIsValid(command shared.Command) error {
	if command == nil {
		return errors.New("command is nil interface")
	}

	if reflect.ValueOf(command).IsNil() {
		return errors.Newf("[%s]: command value is nil pointer", command.CommandName())
	}

	if reflect.ValueOf(command.AggregateID()).IsNil() {
		return errors.Newf("[%s]: command was not properly created", command.CommandName())
	}

	return nil
}
