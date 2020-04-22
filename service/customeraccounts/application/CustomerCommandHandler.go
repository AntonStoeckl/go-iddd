package application

import (
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/shared"
	"github.com/cockroachdb/errors"
)

const maxCustomerCommandHandlerRetries = uint8(10)

type CustomerCommandHandler struct {
	retrieveCustomerEventStream customer.ForRetrievingCustomerEventStreams
	startCustomerEventStream    customer.ForStartingCustomerEventStreams
	appendToCustomerEventStream customer.ForAppendingToCustomerEventStreams
}

func NewCustomerCommandHandler(
	retrieveCustomerEventStream customer.ForRetrievingCustomerEventStreams,
	startCustomerEventStream customer.ForStartingCustomerEventStreams,
	appendToCustomerEventStream customer.ForAppendingToCustomerEventStreams,
) *CustomerCommandHandler {

	return &CustomerCommandHandler{
		retrieveCustomerEventStream: retrieveCustomerEventStream,
		startCustomerEventStream:    startCustomerEventStream,
		appendToCustomerEventStream: appendToCustomerEventStream,
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

	emailAddressValue, err := value.BuildEmailAddress(emailAddress)
	if err != nil {
		return value.CustomerID{}, errors.Wrap(err, wrapWithMsg)
	}

	personNameValue, err := value.BuildPersonName(givenName, familyName)
	if err != nil {
		return value.CustomerID{}, errors.Wrap(err, wrapWithMsg)
	}

	command = domain.BuildRegisterCustomer(
		value.GenerateCustomerID(),
		emailAddressValue,
		value.GenerateConfirmationHash(emailAddressValue.String()),
		personNameValue,
	)

	doRegister := func() error {
		recordedEvents := customer.Register(command)

		if err = h.startCustomerEventStream(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		return nil
	}

	if err = shared.RetryOnConcurrencyConflict(doRegister, maxCustomerCommandHandlerRetries); err != nil {
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

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	confirmationHashValue, err := value.BuildConfirmationHash(confirmationHash)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	command = domain.BuildConfirmCustomerEmailAddress(customerIDValue, confirmationHashValue)

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

	if err := shared.RetryOnConcurrencyConflict(doConfirmEmailAddress, maxCustomerCommandHandlerRetries); err != nil {
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

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	emailAddressValue, err := value.BuildEmailAddress(emailAddress)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	command = domain.BuildChangeCustomerEmailAddress(customerIDValue, emailAddressValue)

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

	if err := shared.RetryOnConcurrencyConflict(doChangeEmailAddress, maxCustomerCommandHandlerRetries); err != nil {
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

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	personNameValue, err := value.BuildPersonName(givenName, familyName)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	command = domain.BuildChangeCustomerName(customerIDValue, personNameValue)

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

	if err := shared.RetryOnConcurrencyConflict(doChangeName, maxCustomerCommandHandlerRetries); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}

func (h *CustomerCommandHandler) DeleteCustomer(customerID string) error {
	var err error
	var command domain.DeleteCustomer
	wrapWithMsg := "customerCommandHandler.DeleteCustomer"

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	command = domain.BuildDeleteCustomer(customerIDValue)

	doDelete := func() error {
		eventStream, err := h.retrieveCustomerEventStream(command.CustomerID())
		if err != nil {
			return err
		}

		recordedEvents := customer.Delete(eventStream, command)

		if err := h.appendToCustomerEventStream(recordedEvents, command.CustomerID()); err != nil {
			return err
		}

		return nil
	}

	if err := shared.RetryOnConcurrencyConflict(doDelete, maxCustomerCommandHandlerRetries); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}
