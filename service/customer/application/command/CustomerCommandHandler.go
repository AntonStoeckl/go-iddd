package command

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/cqrs"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
)

const maxCustomerCommandHandlerRetries = uint8(10)

type CustomerCommandHandler struct {
	customerEvents ForStoringCustomerEvents
}

func NewCustomerCommandHandler(customerEvents ForStoringCustomerEvents) *CustomerCommandHandler {
	return &CustomerCommandHandler{
		customerEvents: customerEvents,
	}
}

func (h *CustomerCommandHandler) RegisterCustomer(command commands.RegisterCustomer) error {
	if err := command.ShouldBeValid(); err != nil {
		return errors.Wrap(err, "customerCommandHandler.RegisterCustomer")
	}

	doRegister := func() error {
		recordedEvents := customer.Register(command)

		if err := h.customerEvents.CreateStreamFrom(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		return nil
	}

	if err := cqrs.RetryCommand(doRegister, maxCustomerCommandHandlerRetries); err != nil {
		return err
	}

	return nil
}

func (h *CustomerCommandHandler) ConfirmCustomerEmailAddress(command commands.ConfirmCustomerEmailAddress) error {
	if err := command.ShouldBeValid(); err != nil {
		return errors.Wrap(err, "customerCommandHandler.ConfirmCustomerEmailAddress")
	}

	doConfirmEmailAddress := func() error {
		eventStream, err := h.customerEvents.EventStreamFor(command.CustomerID())
		if err != nil {
			return err
		}

		if err := h.assertNotDeleted(eventStream); err != nil {
			return errors.Wrap(err, "customerCommandHandler.ChangeCustomerName")
		}

		recordedEvents := customer.ConfirmEmailAddress(eventStream, command)

		if err := h.customerEvents.Add(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		for _, event := range recordedEvents {
			if event.IndicatesAnError() {
				return errors.Mark(errors.New(event.EventName()), lib.ErrDomainConstraintsViolation)
			}
		}

		return nil
	}

	if err := cqrs.RetryCommand(doConfirmEmailAddress, maxCustomerCommandHandlerRetries); err != nil {
		return err
	}

	return nil
}

func (h *CustomerCommandHandler) ChangeCustomerEmailAddress(command commands.ChangeCustomerEmailAddress) error {
	if err := command.ShouldBeValid(); err != nil {
		return errors.Wrap(err, "customerCommandHandler.ChangeCustomerEmailAddress")
	}

	doChangeEmailAddress := func() error {
		eventStream, err := h.customerEvents.EventStreamFor(command.CustomerID())
		if err != nil {
			return err
		}

		if err := h.assertNotDeleted(eventStream); err != nil {
			return errors.Wrap(err, "customerCommandHandler.ChangeCustomerName")
		}

		recordedEvents := customer.ChangeEmailAddress(eventStream, command)

		if err := h.customerEvents.Add(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		return nil
	}

	if err := cqrs.RetryCommand(doChangeEmailAddress, maxCustomerCommandHandlerRetries); err != nil {
		return err
	}

	return nil
}

func (h *CustomerCommandHandler) ChangeCustomerName(command commands.ChangeCustomerName) error {
	if err := command.ShouldBeValid(); err != nil {
		return errors.Wrap(err, "customerCommandHandler.ChangeCustomerName")
	}

	doChangeName := func() error {
		eventStream, err := h.customerEvents.EventStreamFor(command.CustomerID())
		if err != nil {
			return err
		}

		if err := h.assertNotDeleted(eventStream); err != nil {
			return errors.Wrap(err, "customerCommandHandler.ChangeCustomerName")
		}

		recordedEvents := customer.ChangeName(eventStream, command)

		if err := h.customerEvents.Add(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		return nil
	}

	if err := cqrs.RetryCommand(doChangeName, maxCustomerCommandHandlerRetries); err != nil {
		return err
	}

	return nil
}

func (h *CustomerCommandHandler) DeleteCustomer(command commands.DeleteCustomer) error {
	if err := command.ShouldBeValid(); err != nil {
		return errors.Wrap(err, "customerCommandHandler.ChangeCustomerName")
	}

	doChangeName := func() error {
		eventStream, err := h.customerEvents.EventStreamFor(command.CustomerID())
		if err != nil {
			return err
		}

		recordedEvents := customer.Delete(eventStream)

		if err := h.customerEvents.Add(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		return nil
	}

	if err := cqrs.RetryCommand(doChangeName, maxCustomerCommandHandlerRetries); err != nil {
		return err
	}

	return nil
}

func (h *CustomerCommandHandler) assertNotDeleted(eventStream es.DomainEvents) error {
	for _, event := range eventStream {
		switch event.(type) {
		case events.CustomerDeleted:
			return errors.Mark(errors.New("customer not found"), lib.ErrNotFound)
		}
	}

	return nil
}
