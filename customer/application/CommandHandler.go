package application

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"reflect"

	"golang.org/x/xerrors"
)

type CommandHandler struct {
	customers domain.Customers
}

/*** Factory Method ***/

func NewCommandHandler(customers domain.Customers) *CommandHandler {
	return &CommandHandler{customers: customers}
}

/*** Implement shared.CommandHandler ***/

func (handler *CommandHandler) Handle(command shared.Command) error {
	var err error

	if err := handler.assertIsValid(command); err != nil {
		return err
	}

	switch actualCommand := command.(type) {
	case *commands.Register:
		err = handler.register(actualCommand)
	case *commands.ConfirmEmailAddress:
		err = handler.applyToExistingCustomer(actualCommand.ID(), actualCommand)
	default:
		return xerrors.Errorf(
			"commandHandler.Handle: Command [%s] is unknown: %w",
			actualCommand.CommandName(),
			shared.ErrCommandCanNotBeHandled,
		)
	}

	return err
}

/*** Business cases ***/

func (handler *CommandHandler) register(register *commands.Register) error {
	customer := domain.Register(register)

	if err := handler.customers.Save(customer); err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) applyToExistingCustomer(id *values.ID, command shared.Command) error {
	customer, err := handler.customers.FindBy(id)
	if err != nil {
		return err
	}

	if err := customer.Apply(command); err != nil {
		return err
	}

	if err := handler.customers.Save(customer); err != nil {
		return err
	}

	return nil
}

/*** Command Assertions ***/

func (handler *CommandHandler) assertIsValid(command shared.Command) error {
	if command == nil {
		// command is a nil interface
		return xerrors.Errorf(
			"commandHandler.Handle: Command is nil interface: %w",
			shared.ErrCommandIsInvalid,
		)
	}

	if reflect.ValueOf(command).IsNil() {
		// value of the command inteface is a nil pointer
		return xerrors.Errorf(
			"commandHandler.Handle: Command is nil pointer: %w",
			shared.ErrCommandIsInvalid,
		)
	}

	if reflect.ValueOf(command.AggregateIdentifier()).IsNil() {
		// command has no AggregateIdentifier, so it can't have been created with a proper factory method
		return xerrors.Errorf(
			"commandHandler.Handle: Command has no AggregateIdentifier: %w",
			shared.ErrCommandIsInvalid,
		)
	}

	return nil
}
