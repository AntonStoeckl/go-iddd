package application

import (
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
)

const maxCustomerCommandHandlerRetries = uint8(10)

type CustomerCommandHandler struct {
	retrieveCustomerEventStream ForRetrievingCustomerEventStreams
	startCustomerEventStream    ForStartingCustomerEventStreams
	appendToCustomerEventStream ForAppendingToCustomerEventStreams
}

func NewCustomerCommandHandler(
	retrieveCustomerEventStream ForRetrievingCustomerEventStreams,
	startCustomerEventStream ForStartingCustomerEventStreams,
	appendToCustomerEventStream ForAppendingToCustomerEventStreams,
) *CustomerCommandHandler {

	return &CustomerCommandHandler{
		retrieveCustomerEventStream: retrieveCustomerEventStream,
		startCustomerEventStream:    startCustomerEventStream,
		appendToCustomerEventStream: appendToCustomerEventStream,
	}
}

func (h *CustomerCommandHandler) RegisterCustomer(
	customerIDValue value.CustomerID,
	emailAddress string,
	givenName string,
	familyName string,
) error {

	wrapWithMsg := "CustomerCommandHandler.RegisterCustomer"

	emailAddressValue, err := value.BuildUnconfirmedEmailAddress(emailAddress)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	personNameValue, err := value.BuildPersonName(givenName, familyName)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	command := domain.BuildRegisterCustomer(
		customerIDValue,
		emailAddressValue,
		personNameValue,
	)

	doRegister := func() error {
		customerRegistered := customer.Register(command)

		if err := h.startCustomerEventStream(customerRegistered); err != nil {
			return err
		}

		return nil
	}

	if err := shared.RetryOnConcurrencyConflict(doRegister, maxCustomerCommandHandlerRetries); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}

func (h *CustomerCommandHandler) ConfirmCustomerEmailAddress(
	customerID string,
	confirmationHash string,
) error {

	wrapWithMsg := "CustomerCommandHandler.ConfirmCustomerEmailAddress"

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	confirmationHashValue, err := value.BuildConfirmationHash(confirmationHash)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	command := domain.BuildConfirmCustomerEmailAddress(
		customerIDValue,
		confirmationHashValue,
	)

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

	wrapWithMsg := "CustomerCommandHandler.ChangeCustomerEmailAddress"

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	emailAddressValue, err := value.BuildUnconfirmedEmailAddress(emailAddress)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	command := domain.BuildChangeCustomerEmailAddress(
		customerIDValue,
		emailAddressValue,
	)

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

	wrapWithMsg := "CustomerCommandHandler.ChangeCustomerName"

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	personNameValue, err := value.BuildPersonName(givenName, familyName)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	command := domain.BuildChangeCustomerName(
		customerIDValue,
		personNameValue,
	)

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
	wrapWithMsg := "customerCommandHandler.DeleteCustomer"

	customerIDValue, err := value.BuildCustomerID(customerID)
	if err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	command := domain.BuildDeleteCustomer(customerIDValue)

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
