package application

import (
	"errors"
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
)

type commandHandler struct {
	customers domain.Customers
}

func NewCommandHandler(customers domain.Customers) *commandHandler {
	return &commandHandler{customers: customers}
}

func (handler *commandHandler) Handle(command shared.Command) error {
	var err error

	switch command := command.(type) {
	case *commands.Register:
		err = handler.register(command)
	case *commands.ConfirmEmailAddress:
		err = handler.applyToExistingCustomer(command.ID(), command)
	case nil:
		err = errors.New("commandHandler - nil command handled")
	default:
		err = errors.New("commandHandler - unknown command handled")
	}

	return err
}

func (handler *commandHandler) register(register *commands.Register) error {
	customer := handler.customers.New()

	if err := customer.Apply(register); err != nil {
		return err
	}

	if err := handler.customers.Save(customer); err != nil {
		return err
	}

	return nil
}

func (handler *commandHandler) applyToExistingCustomer(id *values.ID, command shared.Command) error {
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
