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

type CustomersWithPersistance interface {
	domain.Customers
	shared.PersistsEventsourcedAggregates
}

type CommandHandler struct {
	customers CustomersWithPersistance
}

/*** Factory Method ***/

func NewCommandHandler(customers CustomersWithPersistance) *CommandHandler {
	return &CommandHandler{customers: customers}
}

/*** Implement shared.CommandHandler ***/

func (handler *CommandHandler) Handle(command shared.Command) error {
	var err error

	if err := handler.assertIsValid(command); err != nil {
		return xerrors.Errorf(
			"commandHandler.Handle: %s: %w",
			err,
			shared.ErrCommandIsInvalid,
		)
	}

	switch actualCommand := command.(type) {
	case *commands.Register:
		err = handler.register(actualCommand)
	case *commands.ConfirmEmailAddress:
		err = handler.confirmEmailAddress(actualCommand)
	case *commands.ChangeEmailAddress:
		err = handler.changeEmailAddress(actualCommand)
	default:
		return xerrors.Errorf(
			"commandHandler.Handle: Command [%s] is unknown: %w",
			actualCommand.CommandName(),
			shared.ErrCommandCanNotBeHandled,
		)
	}

	if err != nil {
		return xerrors.Errorf(
			"commandHandler.Handle: %s: %w",
			command.CommandName(),
			err,
		)
	}

	return nil
}

/*** Business cases ***/

func (handler *CommandHandler) register(register *commands.Register) error {
	newCustomer := domain.NewCustomerWith(register)

	if err := handler.customers.Register(newCustomer); err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) confirmEmailAddress(confirmEmailAddress *commands.ConfirmEmailAddress) error {
	err := handler.applyToCustomer(
		confirmEmailAddress.ID(),
		func(customer domain.Customer) error {
			return customer.Execute(confirmEmailAddress)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) changeEmailAddress(changeEmailAddress *commands.ChangeEmailAddress) error {
	err := handler.applyToCustomer(
		changeEmailAddress.ID(),
		func(customer domain.Customer) error {
			return customer.Execute(changeEmailAddress)
		},
	)

	if err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) applyToCustomer(id *values.ID, callback func(customer domain.Customer) error) error {
	customer, err := handler.customers.Of(id)
	if err != nil {
		return err
	}

	if err := callback(customer); err != nil {
		return err
	}

	if err := handler.customers.Persist(customer); err != nil {
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

	if reflect.ValueOf(command.AggregateID()).IsNil() {
		return fmt.Errorf("[%s]: command was not properly created", command.CommandName())
	}

	return nil
}
