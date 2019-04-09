package application

import (
	"errors"
	"go-iddd/customer/model"
	"go-iddd/customer/model/commands"
	"go-iddd/shared"
)

type commandHandler struct {
	customers model.Customers
}

func NewCommandHandler(customers model.Customers) *commandHandler {
	return &commandHandler{customers: customers}
}

func (handler *commandHandler) Handle(command shared.Command) error {
	var err error

	switch command := command.(type) {
	case commands.Register:
		err = handler.register(command)
	case commands.ConfirmEmailAddress:
		err = handler.confirmEmailAddress(command)
	case nil:
		err = errors.New("commandHandler - nil command handled")
	default:
		err = errors.New("commandHandler - unknown command handled")
	}

	return err
}

func (handler *commandHandler) register(register commands.Register) error {
	newCustomer := model.NewUnregisteredCustomer()

	if err := newCustomer.Apply(register); err != nil {
		return err
	}

	if err := handler.customers.Save(newCustomer); err != nil {
		return err
	}

	return nil
}

func (handler *commandHandler) confirmEmailAddress(confirmEmailAddress commands.ConfirmEmailAddress) error {
	customer, err := handler.customers.FindBy(confirmEmailAddress.ID())
	if err != nil {
		return err
	}

	if err := customer.Apply(confirmEmailAddress); err != nil {
		return err
	}

	if err := handler.customers.Save(customer); err != nil {
		return err
	}

	return nil
}
