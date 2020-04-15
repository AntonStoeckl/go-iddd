package application

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer/value"
	"github.com/cockroachdb/errors"
)

const maxCustomerCommandHandlerRetries = uint8(10)

type CustomerCommandHandler struct {
	retrieveCustomerEventStream ForRetrievingCustomerEventStreams
	startCustomerEventStream    ForStartingCustomerEventStreams
	appendToCustomerEventStream ForAppendingToCustomerEventStreams
	retryCommand                ForRetryingCommands
}

func NewCustomerCommandHandler(
	retrieveCustomerEventStream ForRetrievingCustomerEventStreams,
	startCustomerEventStream ForStartingCustomerEventStreams,
	appendToCustomerEventStream ForAppendingToCustomerEventStreams,
	retryCommand ForRetryingCommands,
) *CustomerCommandHandler {

	return &CustomerCommandHandler{
		retrieveCustomerEventStream: retrieveCustomerEventStream,
		startCustomerEventStream:    startCustomerEventStream,
		appendToCustomerEventStream: appendToCustomerEventStream,
		retryCommand:                retryCommand,
	}
}

func (h *CustomerCommandHandler) RegisterCustomer(
	emailAddress string,
	givenName string,
	familyName string,
) (value.CustomerID, error) {

	var err error
	var command domain.RegisterCustomer
	wrapWithMsg := "customerCommandHandler.RegisterCustomer"

	if command, err = domain.BuildRegisterCustomer(emailAddress, givenName, familyName); err != nil {
		return value.CustomerID{}, errors.Wrap(err, wrapWithMsg)
	}

	doRegister := func() error {
		recordedEvents := customer.Register(command)

		if err = h.startCustomerEventStream(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		return nil
	}

	if err = h.retryCommand(doRegister, maxCustomerCommandHandlerRetries); err != nil {
		return value.CustomerID{}, errors.Wrap(err, wrapWithMsg)
	}

	return command.CustomerID(), nil
}

func (h *CustomerCommandHandler) ConfirmCustomerEmailAddress(
	customerID string,
	confirmationHash string,
) error {

	var err error
	var command domain.ConfirmCustomerEmailAddress
	wrapWithMsg := "customerCommandHandler.ConfirmCustomerEmailAddress"

	if command, err = domain.BuildConfirmCustomerEmailAddress(customerID, confirmationHash); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	doConfirmEmailAddress := func() error {
		eventStream, err := h.retrieveCustomerEventStream(command.CustomerID())
		if err != nil {
			return err
		}

		recordedEvents, err := customer.ConfirmEmailAddress(eventStream, command)
		if err != nil {
			return err
		}

		if err := h.appendToCustomerEventStream(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		for _, event := range recordedEvents {
			if isError := event.IsFailureEvent(); isError {
				return event.FailureReason()
			}
		}

		return nil
	}

	if err := h.retryCommand(doConfirmEmailAddress, maxCustomerCommandHandlerRetries); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}

func (h *CustomerCommandHandler) ChangeCustomerEmailAddress(
	customerID string,
	emailAddress string,
) error {

	var err error
	var command domain.ChangeCustomerEmailAddress
	wrapWithMsg := "customerCommandHandler.ChangeCustomerEmailAddress"

	if command, err = domain.BuildChangeCustomerEmailAddress(customerID, emailAddress); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	doChangeEmailAddress := func() error {
		eventStream, err := h.retrieveCustomerEventStream(command.CustomerID())
		if err != nil {
			return err
		}

		recordedEvents, err := customer.ChangeEmailAddress(eventStream, command)
		if err != nil {
			return err
		}

		if err := h.appendToCustomerEventStream(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		return nil
	}

	if err := h.retryCommand(doChangeEmailAddress, maxCustomerCommandHandlerRetries); err != nil {
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
	var command domain.ChangeCustomerName
	wrapWithMsg := "customerCommandHandler.ChangeCustomerName"

	if command, err = domain.BuildChangeCustomerName(customerID, givenName, familyName); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	doChangeName := func() error {
		eventStream, err := h.retrieveCustomerEventStream(command.CustomerID())
		if err != nil {
			return err
		}

		recordedEvents, err := customer.ChangeName(eventStream, command)
		if err != nil {
			return err
		}

		if err := h.appendToCustomerEventStream(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		return nil
	}

	if err := h.retryCommand(doChangeName, maxCustomerCommandHandlerRetries); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}

func (h *CustomerCommandHandler) DeleteCustomer(customerID string) error {
	var err error
	var command domain.DeleteCustomer
	wrapWithMsg := "customerCommandHandler.DeleteCustomer"

	if command, err = domain.BuildDeleteCustomer(customerID); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	doDelete := func() error {
		eventStream, err := h.retrieveCustomerEventStream(command.CustomerID())
		if err != nil {
			return err
		}

		recordedEvents := customer.Delete(eventStream)

		if err := h.appendToCustomerEventStream(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		return nil
	}

	if err := h.retryCommand(doDelete, maxCustomerCommandHandlerRetries); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}
