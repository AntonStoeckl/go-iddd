package application

import (
	"errors"
	"go-iddd/customer/model"
	"go-iddd/customer/model/commands"
	"go-iddd/shared"
	"reflect"
)

type commandHandlerWithReflection struct {
	customers model.Customers
}

func NewCommandHandlerWithReflection(persons model.Customers) *commandHandlerWithReflection {
	return &commandHandlerWithReflection{customers: persons}
}

func (handler *commandHandlerWithReflection) Handle(command shared.Command) error {
	if command == nil {
		return errors.New("commandHandler - nil command handled")
	}

	method, ok := reflect.TypeOf(handler).MethodByName(command.CommandName())
	if !ok {
		return errors.New("commandHandler - unknown command handled")
	}

	in := make([]reflect.Value, 2)
	in[0] = reflect.ValueOf(handler) // method receiver
	in[1] = reflect.ValueOf(command) // first input param - the command

	response := method.Func.Call(in)

	switch response[0].Interface().(type) {
	case error:
		return response[0].Interface().(error)
	case nil:
		return nil
	default:
		return errors.New("commandHandler - unexpected type returned when command was handled")
	}
}

func (handler *commandHandlerWithReflection) Register(register commands.Register) error {
	newCustomer := model.NewUnregisteredCustomer()

	if err := newCustomer.Apply(register); err != nil {
		return err
	}

	if err := handler.customers.Save(newCustomer); err != nil {
		return err
	}

	return nil
}

func (handler *commandHandlerWithReflection) ConfirmEmailAddress(confirmEmailAddress commands.ConfirmEmailAddress) error {
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
