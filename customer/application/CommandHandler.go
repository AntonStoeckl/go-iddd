package application

import (
	"errors"
	"fmt"
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
		return xerrors.Errorf("commandHandler.Handle: %s: %w", err, shared.ErrCommandIsInvalid)
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
