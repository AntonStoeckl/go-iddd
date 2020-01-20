package application

import (
	"database/sql"
	"go-iddd/service/customer/domain"
	"go-iddd/service/customer/domain/commands"
	"go-iddd/service/customer/domain/events"
	"go-iddd/service/lib"

	"github.com/cockroachdb/errors"
)

const maxCommandHandlerRetries = uint8(10)

type CommandHandler struct {
	sessionStarter StartsCustomersSession
	db             *sql.DB
}

func NewCommandHandler(sessionStarter StartsCustomersSession, db *sql.DB) *CommandHandler {
	return &CommandHandler{
		sessionStarter: sessionStarter,
		db:             db,
	}
}

func (handler *CommandHandler) Register(command commands.Register) error {
	if err := command.ShouldBeValid(); err != nil {
		return errors.Wrap(err, "commandHandler.Register")
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

func (handler *CommandHandler) handleRetry(command lib.Command) error {
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

func (handler *CommandHandler) handleSession(command lib.Command) error {
	tx, errTx := handler.db.Begin()
	if errTx != nil {
		return errors.Mark(errTx, lib.ErrTechnical)
	}

	customers := handler.sessionStarter.StartSession(tx)

	// call next method in chain
	err := handler.handleCommand(customers, command)

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
	customers Customers,
	command lib.Command,
) error {

	var err error

	switch actualCommand := command.(type) {
	case commands.Register:
		err = handler.register(customers, actualCommand)
	case commands.ConfirmEmailAddress:
		err = handler.confirmEmailAddress(customers, actualCommand)
	case commands.ChangeEmailAddress:
		err = handler.changeEmailAddress(customers, actualCommand)
	}

	if err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) register(
	customers Customers,
	register commands.Register,
) error {

	recordedEvents := domain.RegisterCustomer(register)

	if err := customers.Register(register.CustomerID(), recordedEvents); err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) confirmEmailAddress(
	customers Customers,
	confirmEmailAddress commands.ConfirmEmailAddress,
) error {

	eventStream, err := customers.EventStream(confirmEmailAddress.CustomerID())
	if err != nil {
		return err
	}

	recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)

	if err := customers.Persist(confirmEmailAddress.CustomerID(), recordedEvents); err != nil {
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
	customers Customers,
	changeEmailAddress commands.ChangeEmailAddress,
) error {

	eventStream, err := customers.EventStream(changeEmailAddress.CustomerID())
	if err != nil {
		return err
	}

	recordedEvents := domain.ChangeEmailAddress(eventStream, changeEmailAddress)

	if err := customers.Persist(changeEmailAddress.CustomerID(), recordedEvents); err != nil {
		return err
	}

	return nil
}
