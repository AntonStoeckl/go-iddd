package application

import (
    "errors"
    "go-iddd/customer/domain"
    "go-iddd/customer/domain/valueobjects"
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
    case domain.Register:
        err = handler.register(command)
    case domain.ConfirmEmailAddress:
        err = handler.applyToExistingCustomer(command.ID(), command)
    case nil:
        err = errors.New("commandHandler - nil command handled")
    default:
        err = errors.New("commandHandler - unknown command handled")
    }

    return err
}

func (handler *commandHandler) register(register domain.Register) error {
    customer := handler.customers.New()

    if err := customer.Apply(register); err != nil {
        return err
    }

    if err := handler.customers.Save(customer); err != nil {
        return err
    }

    return nil
}

func (handler *commandHandler) applyToExistingCustomer(id valueobjects.ID, command shared.Command) error {
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
