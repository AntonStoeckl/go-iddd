package command

import (
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/customer"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/lib"
	"go-iddd/service/lib/cqrs"

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

		recordedEvents := customer.ConfirmEmailAddress(eventStream, command)

		if err := h.customerEvents.Add(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		for _, event := range recordedEvents {
			switch actualEvent := event.(type) {
			case events.CustomerEmailAddressConfirmationFailed:
				return errors.Mark(errors.New(actualEvent.EventName()), lib.ErrDomainConstraintsViolation)
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
