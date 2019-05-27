package domain

import (
	"errors"
	"fmt"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"reflect"
	"strings"

	"golang.org/x/xerrors"
)

type Customer interface {
	Apply(cmd shared.Command) error

	shared.Aggregate
	shared.RecordsEvents
}

type customer struct {
	id                      *values.ID
	confirmableEmailAddress *values.ConfirmableEmailAddress
	personName              *values.PersonName
	recordedEvents          shared.EventStream
}

func blankCustomer() *customer {
	return &customer{}
}

/*** Implement Customer ***/

func (customer *customer) Apply(command shared.Command) error {
	var err error

	if err := customer.assertIsValid(command); err != nil {
		return xerrors.Errorf("customer.Apply: %s: %w", err, shared.ErrCommandIsInvalid)
	}

	switch actualCommand := command.(type) {
	case *commands.ConfirmEmailAddress:
		err = customer.confirmEmailAddress(actualCommand)
	case *commands.Register:
		return xerrors.Errorf("customer.Apply: customer is already registered: %w", shared.ErrCommandCanNotBeHandled)
	default:
		return xerrors.Errorf("customer.Apply: [%s]: command is unknown: %w", command.CommandName(), shared.ErrCommandCanNotBeHandled)
	}

	if err != nil {
		return err
	}

	return nil
}

/*** Implement shared.Aggregate ****/

func (customer *customer) AggregateIdentifier() shared.AggregateIdentifier {
	return customer.id
}

func (customer *customer) AggregateName() string {
	aggregateType := reflect.TypeOf(customer).String()
	aggregateTypeParts := strings.Split(aggregateType, ".")
	aggregateName := aggregateTypeParts[len(aggregateTypeParts)-1]

	return strings.Title(aggregateName)
}

/*** EventSourcing ***/

func ReconstituteCustomerFrom(eventStream shared.EventStream) (Customer, error) {
	newCustomer := blankCustomer()

	if !eventStream.FirstEventIsOfSameTypeAs(&events.Registered{}) {
		return nil, xerrors.Errorf("ReconstituteCustomerFrom: %w", shared.ErrInvalidEventStream)
	}

	for _, event := range eventStream {
		newCustomer.when(event)
	}

	return newCustomer, nil
}

func (customer *customer) recordThat(event shared.DomainEvent) {
	customer.recordedEvents = append(customer.recordedEvents, event)
	customer.when(event)
}

func (customer *customer) when(event shared.DomainEvent) {
	switch actualEvent := event.(type) {
	case *events.Registered:
		customer.whenItWasRegistered(actualEvent)
	case *events.EmailAddressConfirmed:
		customer.whenEmailAddressWasConfirmed(actualEvent)
	}
}

/*** Implement shared.RecordsEvents ****/

func (customer *customer) RecordedEvents() shared.EventStream {
	currentEvents := customer.recordedEvents
	customer.recordedEvents = nil

	return currentEvents
}

/*** Command Assertions ***/

func (customer *customer) assertIsValid(command shared.Command) error {
	if command == nil {
		return errors.New("command is nil interface")
	}

	if reflect.ValueOf(command).IsNil() {
		return fmt.Errorf("[%s]: command value is nil pointer", command.CommandName())
	}

	if reflect.ValueOf(command.AggregateIdentifier()).IsNil() {
		return fmt.Errorf("[%s]: command was not properly created", command.CommandName())
	}

	return nil
}
