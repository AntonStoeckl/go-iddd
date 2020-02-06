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

func (handler *CommandHandler) Register(register commands.Register) error {
	if err := register.ShouldBeValid(); err != nil {
		return errors.Wrap(err, "commandHandler.CreateStreamFrom")
	}

	if err := cqrs.RetryCommand(
		func() error {
			recordedEvents := domain.RegisterCustomer(register)

			if err := handler.customerEvents.CreateStreamFrom(recordedEvents, register.CustomerID()); err != nil {
				return err
			}

			return nil
		},
		maxCommandHandlerRetries,
	); err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) ConfirmEmailAddress(confirmEmailAddress commands.ConfirmEmailAddress) error {
	if err := confirmEmailAddress.ShouldBeValid(); err != nil {
		return errors.Wrap(err, "commandHandler.ConfirmEmailAddress")
	}

	if err := cqrs.RetryCommand(
		func() error {
			eventStream, err := handler.customerEvents.EventStreamFor(confirmEmailAddress.CustomerID())
			if err != nil {
				return err
			}

			recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)

			if err := handler.customerEvents.Add(recordedEvents, confirmEmailAddress.CustomerID()); err != nil {
				return err
			}

			for _, event := range recordedEvents {
				switch actualEvent := event.(type) {
				case events.EmailAddressConfirmationFailed:
					return errors.Mark(errors.New(actualEvent.EventName()), lib.ErrDomainConstraintsViolation)
				}
			}

			return nil
		},
		maxCommandHandlerRetries,
	); err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) ChangeEmailAddress(changeEmailAddress commands.ChangeEmailAddress) error {
	if err := changeEmailAddress.ShouldBeValid(); err != nil {
		return errors.Wrap(err, "commandHandler.ChangeEmailAddress")
	}

	if err := cqrs.RetryCommand(
		func() error {
			eventStream, err := handler.customerEvents.EventStreamFor(changeEmailAddress.CustomerID())
			if err != nil {
				return err
			}

			recordedEvents := domain.ChangeEmailAddress(eventStream, changeEmailAddress)

			if err := handler.customerEvents.Add(recordedEvents, changeEmailAddress.CustomerID()); err != nil {
				return err
			}

			return nil
		},
		maxCommandHandlerRetries,
	); err != nil {
		return err
	}

	return nil
}
