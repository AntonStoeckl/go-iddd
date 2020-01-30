package application

import (
	"database/sql"
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
	db             *sql.DB
}

func NewCommandHandler(customerEvents ForStoringCustomerEvents, db *sql.DB) *CommandHandler {
	return &CommandHandler{
		customerEvents: customerEvents,
		db:             db,
	}
}

func (handler *CommandHandler) Register(command commands.Register) error {
	if err := command.ShouldBeValid(); err != nil {
		return errors.Wrap(err, "commandHandler.CreateStreamFrom")
	}

	if err := handler.handleRetry(command); err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) ConfirmEmailAddress(command commands.ConfirmEmailAddress) error {
	if err := command.ShouldBeValid(); err != nil {
		return errors.Wrap(err, "commandHandler.ConfirmEmailAddress")
	}

	if err := handler.handleRetry(command); err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) ChangeEmailAddress(command commands.ChangeEmailAddress) error {
	if err := command.ShouldBeValid(); err != nil {
		return errors.Wrap(err, "commandHandler.ChangeEmailAddress")
	}

	if err := handler.handleRetry(command); err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) handleRetry(command cqrs.Command) error {
	var err error
	var retries uint8

	for retries = 0; retries < maxCommandHandlerRetries; retries++ {
		// call next method in chain
		if err = handler.handleSession(command); err == nil {
			break // no need to retry, handling was successful
		}

		if errors.Is(err, lib.ErrConcurrencyConflict) {
			continue // retry to resolve the concurrency conflict
		} else {
			break // don't retry for different errors
		}
	}

	if err != nil {
		if retries == maxCommandHandlerRetries {
			return errors.Wrap(err, lib.ErrMaxRetriesExceeded.Error())
		}

		return err // either to many retries or a different error
	}

	return nil
}

func (handler *CommandHandler) handleSession(command cqrs.Command) error {
	tx, errTx := handler.db.Begin()
	if errTx != nil {
		return errors.Mark(errTx, lib.ErrTechnical)
	}

	// call next method in chain
	err := handler.handleCommand(command, tx)

	if err != nil {
		if !errors.Is(err, lib.ErrDomainConstraintsViolation) {
			_ = tx.Rollback()

			return err
		}
	}

	if errTx := tx.Commit(); errTx != nil {
		return errors.Mark(errTx, lib.ErrTechnical)
	}

	if err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) handleCommand(
	command cqrs.Command,
	tx *sql.Tx,
) error {

	var err error

	switch actualCommand := command.(type) {
	case commands.Register:
		err = handler.register(actualCommand, tx)
	case commands.ConfirmEmailAddress:
		err = handler.confirmEmailAddress(actualCommand, tx)
	case commands.ChangeEmailAddress:
		err = handler.changeEmailAddress(actualCommand, tx)
	}

	if err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) register(
	register commands.Register,
	tx *sql.Tx,
) error {

	recordedEvents := domain.RegisterCustomer(register)

	if err := handler.customerEvents.CreateStreamFrom(recordedEvents, register.CustomerID(), tx); err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) confirmEmailAddress(
	confirmEmailAddress commands.ConfirmEmailAddress,
	tx *sql.Tx,
) error {

	eventStream, err := handler.customerEvents.EventStreamFor(confirmEmailAddress.CustomerID())
	if err != nil {
		return err
	}

	recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)

	if err := handler.customerEvents.Add(recordedEvents, confirmEmailAddress.CustomerID(), tx); err != nil {
		return err
	}

	for _, event := range recordedEvents {
		switch actualEvent := event.(type) {
		case events.EmailAddressConfirmationFailed:
			return errors.Mark(errors.New(actualEvent.EventName()), lib.ErrDomainConstraintsViolation)
		}
	}

	return nil
}

func (handler *CommandHandler) changeEmailAddress(
	changeEmailAddress commands.ChangeEmailAddress,
	tx *sql.Tx,
) error {

	eventStream, err := handler.customerEvents.EventStreamFor(changeEmailAddress.CustomerID())
	if err != nil {
		return err
	}

	recordedEvents := domain.ChangeEmailAddress(eventStream, changeEmailAddress)

	if err := handler.customerEvents.Add(recordedEvents, changeEmailAddress.CustomerID(), tx); err != nil {
		return err
	}

	return nil
}
