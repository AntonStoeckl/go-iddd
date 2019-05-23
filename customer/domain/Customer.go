package domain

import (
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"reflect"

	"golang.org/x/xerrors"
)

type Customer interface {
	Apply(cmd shared.Command) error

	shared.Aggregate
}

type customer struct {
	id                      *values.ID
	confirmableEmailAddress *values.ConfirmableEmailAddress
	personName              *values.PersonName
}

func blankCustomer() *customer {
	return &customer{}
}

/***** Factory methods *****/

func Register(given *commands.Register) Customer {
	newCustomer := blankCustomer()

	newCustomer.id = given.ID()
	newCustomer.confirmableEmailAddress = given.EmailAddress().ToConfirmable()
	newCustomer.personName = given.PersonName()

	return newCustomer
}

/***** Implement Customer (own methods) *****/

func (customer *customer) Apply(command shared.Command) error {
	var err error

	if err := customer.assertIsValid(command); err != nil {
		return err
	}

	switch actualCommand := command.(type) {
	case *commands.ConfirmEmailAddress:
		err = customer.confirmEmailAddress(actualCommand)
	case *commands.Register:
		return customer.alreadyRegisteredError()
	default:
		return customer.unknownCommandError(command.CommandName())
	}

	if err != nil {
		return err
	}

	return nil
}

/***** Customer business cases (other than Register) *****/

func (customer *customer) confirmEmailAddress(given *commands.ConfirmEmailAddress) error {
	if customer.confirmableEmailAddress.IsConfirmed() {
		return nil
	}

	confirmableEmailAddress, err := customer.confirmableEmailAddress.Confirm(
		given.EmailAddress(),
		given.ConfirmationHash(),
	)

	if err != nil {
		return xerrors.Errorf(
			"customer.confirmEmailAddress -> %s: %w",
			err,
			shared.ErrDomainConstraintsViolation,
		)
	}

	customer.confirmableEmailAddress = confirmableEmailAddress

	return nil
}

/***** Implement shared.Aggregate ****/

func (customer *customer) AggregateIdentifier() shared.AggregateIdentifier {
	return customer.id
}

func (customer *customer) AggregateName() string {
	return shared.BuildAggregateNameFor(customer)
}

/***** Command Assertions *****/

func (customer *customer) assertIsValid(command shared.Command) error {
	if command == nil {
		return customer.commandIsNilInterfaceError()
	}

	if reflect.ValueOf(command).IsNil() {
		return customer.commandValueIsNilPointerError()
	}

	if reflect.ValueOf(command.AggregateIdentifier()).IsNil() {
		return customer.commandWasNotProperlyCreatedError()
	}

	return nil
}

/***** Wrapped Errors *****/

func (customer *customer) alreadyRegisteredError() error {
	return xerrors.Errorf(
		"customer.Apply: customer is already registered: %w",
		shared.ErrCommandCanNotBeHandled,
	)
}

func (customer *customer) unknownCommandError(commandName string) error {
	return xerrors.Errorf(
		"customer.Apply: [%s]: command is unknown: %w",
		commandName,
		shared.ErrCommandCanNotBeHandled,
	)
}

func (customer *customer) commandIsNilInterfaceError() error {
	return xerrors.Errorf(
		"commandHandler.Handle: Command is nil interface: %w",
		shared.ErrCommandIsInvalid,
	)
}

func (customer *customer) commandValueIsNilPointerError() error {
	return xerrors.Errorf(
		"commandHandler.Handle: Command value is nil pointer: %w",
		shared.ErrCommandIsInvalid,
	)
}

func (customer *customer) commandWasNotProperlyCreatedError() error {
	return xerrors.Errorf(
		"commandHandler.Handle: Command was not properly created: %w",
		shared.ErrCommandIsInvalid,
	)
}
