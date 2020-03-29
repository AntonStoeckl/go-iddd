package command

import (
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/cqrs"
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

func (h *CustomerCommandHandler) RegisterCustomer(
	emailAddress string,
	givenName string,
	familyName string,
) (values.CustomerID, error) {

	var err error
	var command commands.RegisterCustomer
	wrapWithMsg := "customerCommandHandler.RegisterCustomer"

	if command, err = commands.BuildRegisterCustomer(emailAddress, givenName, familyName); err != nil {
		return values.CustomerID{}, errors.Wrap(err, wrapWithMsg)
	}

	doRegister := func() error {
		recordedEvents := customer.Register(command)

		if err = h.customerEvents.CreateStreamFrom(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		return nil
	}

	if err = cqrs.RetryCommand(doRegister, maxCustomerCommandHandlerRetries); err != nil {
		return values.CustomerID{}, errors.Wrap(err, wrapWithMsg)
	}

	return command.CustomerID(), nil
}

func (h *CustomerCommandHandler) ConfirmCustomerEmailAddress(
	customerID string,
	confirmationHash string,
) error {

	var err error
	var command commands.ConfirmCustomerEmailAddress
	wrapWithMsg := "customerCommandHandler.ConfirmCustomerEmailAddress"

	if command, err = commands.BuildConfirmCustomerEmailAddress(customerID, confirmationHash); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	doConfirmEmailAddress := func() error {
		eventStream, err := h.customerEvents.EventStreamFor(command.CustomerID())
		if err != nil {
			return err
		}

		recordedEvents, err := customer.ConfirmEmailAddress(eventStream, command)
		if err != nil {
			return err
		}

		if err := h.customerEvents.Add(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		for _, event := range recordedEvents {
			if isError, reason := event.IndicatesAnError(); isError {
				return errors.Mark(errors.New(reason), lib.ErrDomainConstraintsViolation)
			}
		}

		return nil
	}

	if err := cqrs.RetryCommand(doConfirmEmailAddress, maxCustomerCommandHandlerRetries); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}

func (h *CustomerCommandHandler) ChangeCustomerEmailAddress(
	customerID string,
	emailAddress string,
) error {

	var err error
	var command commands.ChangeCustomerEmailAddress
	wrapWithMsg := "customerCommandHandler.ChangeCustomerEmailAddress"

	if command, err = commands.BuildChangeCustomerEmailAddress(customerID, emailAddress); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	doChangeEmailAddress := func() error {
		eventStream, err := h.customerEvents.EventStreamFor(command.CustomerID())
		if err != nil {
			return err
		}

		recordedEvents, err := customer.ChangeEmailAddress(eventStream, command)
		if err != nil {
			return err
		}

		if err := h.customerEvents.Add(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		return nil
	}

	if err := cqrs.RetryCommand(doChangeEmailAddress, maxCustomerCommandHandlerRetries); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}

func (h *CustomerCommandHandler) ChangeCustomerName(
	customerID string,
	givenName string,
	familyName string,
) error {

	var err error
	var command commands.ChangeCustomerName
	wrapWithMsg := "customerCommandHandler.ChangeCustomerName"

	if command, err = commands.BuildChangeCustomerName(customerID, givenName, familyName); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	doChangeName := func() error {
		eventStream, err := h.customerEvents.EventStreamFor(command.CustomerID())
		if err != nil {
			return err
		}

		recordedEvents, err := customer.ChangeName(eventStream, command)
		if err != nil {
			return err
		}

		if err := h.customerEvents.Add(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		return nil
	}

	if err := cqrs.RetryCommand(doChangeName, maxCustomerCommandHandlerRetries); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}

func (h *CustomerCommandHandler) DeleteCustomer(customerID string) error {
	var err error
	var command commands.DeleteCustomer
	wrapWithMsg := "customerCommandHandler.DeleteCustomer"

	if command, err = commands.BuildDeleteCustomer(customerID); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	doDelete := func() error {
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

	if err := cqrs.RetryCommand(doDelete, maxCustomerCommandHandlerRetries); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}
