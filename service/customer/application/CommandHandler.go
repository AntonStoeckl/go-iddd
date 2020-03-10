package application

import (
	"go-iddd/service/customer/application/domain"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/lib"
	"go-iddd/service/lib/cqrs"

	"github.com/cockroachdb/errors"
)

const maxCommandHandlerRetries = uint8(10)

type CommandHandler struct {
	customerEvents ForStoringCustomerEvents
}

func NewCommandHandler(customerEvents ForStoringCustomerEvents) *CommandHandler {
	return &CommandHandler{
		customerEvents: customerEvents,
	}
}

func (handler *CommandHandler) RegisterCustomer(command commands.RegisterCustomer) error {
	if err := command.ShouldBeValid(); err != nil {
		return errors.Wrap(err, "commandHandler.RegisterCustomer")
	}

	doRegister := func() error {
		recordedEvents := domain.RegisterCustomer(command)

		if err := handler.customerEvents.CreateStreamFrom(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		return nil
	}

	if err := cqrs.RetryCommand(doRegister, maxCommandHandlerRetries); err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) ConfirmCustomerEmailAddress(command commands.ConfirmCustomerEmailAddress) error {
	if err := command.ShouldBeValid(); err != nil {
		return errors.Wrap(err, "commandHandler.ConfirmCustomerEmailAddress")
	}

	doConfirmEmailAddress := func() error {
		eventStream, err := handler.customerEvents.EventStreamFor(command.CustomerID())
		if err != nil {
			return err
		}

		recordedEvents := domain.ConfirmCustomerEmailAddress(eventStream, command)

		if err := handler.customerEvents.Add(recordedEvents, command.CustomerID()); err != nil {
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

	if err := cqrs.RetryCommand(doConfirmEmailAddress, maxCommandHandlerRetries); err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) ChangeCustomerEmailAddress(command commands.ChangeCustomerEmailAddress) error {
	if err := command.ShouldBeValid(); err != nil {
		return errors.Wrap(err, "commandHandler.ChangeCustomerEmailAddress")
	}

	doChangeEmailAddress := func() error {
		eventStream, err := handler.customerEvents.EventStreamFor(command.CustomerID())
		if err != nil {
			return err
		}

		recordedEvents := domain.ChangeCustomerEmailAddress(eventStream, command)

		if err := handler.customerEvents.Add(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		return nil
	}

	if err := cqrs.RetryCommand(doChangeEmailAddress, maxCommandHandlerRetries); err != nil {
		return err
	}

	return nil
}
