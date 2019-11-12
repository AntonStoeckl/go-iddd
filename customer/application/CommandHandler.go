package application

import (
	"database/sql"
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/shared"
	"reflect"

	"github.com/cockroachdb/errors"
)

const maxCommandHandlerRetries = 10

type CommandHandler struct {
	sessionStarter StartsCustomersSession
	db             *sql.DB
}

/*** Factory Method ***/

func NewCommandHandler(sessionStarter StartsCustomersSession, db *sql.DB) *CommandHandler {
	return &CommandHandler{
		sessionStarter: sessionStarter,
		db:             db,
	}
}

/*** Implement shared.CommandHandler ***/

func (handler *CommandHandler) Handle(command shared.Command) error {
	if err := handler.assertIsValid(command); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrCommandIsInvalid), "commandHandler.Handle")
	}

	if err := handler.assertIsKnown(command); err != nil {
		return errors.Wrap(errors.Mark(err, shared.ErrCommandIsUnknown), "commandHandler.Handle")
	}

	if err := handler.handleRetry(command); err != nil {
		return errors.Wrapf(err, "commandHandler.Handle: [%s]", command.CommandName())
	}

	return nil
}

/*** Chain of handler functions ***/

func (handler *CommandHandler) handleRetry(command shared.Command) error {
	var err error
	var retries uint

	for retries = 0; retries < maxCommandHandlerRetries; retries++ {
		// call next method in chain
		if err = handler.handleSession(command); err == nil {
			break // no need to retry, handling was successful
		}

		if errors.Is(err, shared.ErrConcurrencyConflict) {
			continue // retry to resolve the concurrency conflict
		} else {
			break // don't retry for different errors
		}
	}

	if err != nil {
		if retries == maxCommandHandlerRetries {
			return errors.Wrap(err, shared.ErrMaxRetriesExceeded.Error())
		}

		return err // either to many retries or a different error
	}

	return nil
}

func (handler *CommandHandler) handleSession(command shared.Command) error {
	tx, errTx := handler.db.Begin()
	if errTx != nil {
		return errors.Mark(errTx, shared.ErrTechnical)
	}

	customers := handler.sessionStarter.StartSession(tx)

	// call next method in chain
	err := handler.handleCommand(customers, command)

	if err != nil && !errors.Is(err, shared.ErrDomainConstraintsViolation) {
		_ = tx.Rollback()

		return err
	}

	if errTx := tx.Commit(); errTx != nil {
		return errors.Mark(errTx, shared.ErrTechnical)
	}

	if err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) handleCommand(
	customers Customers,
	command shared.Command,
) error {

	var err error

	switch actualCommand := command.(type) {
	case *commands.Register:
		err = handler.register(customers, actualCommand)
	case *commands.ConfirmEmailAddress:
		err = handler.confirmEmailAddress(customers, actualCommand)
	case *commands.ChangeEmailAddress:
		err = handler.changeEmailAddress(customers, actualCommand)
	}

	if err != nil {
		return err
	}

	return nil
}

/*** Business cases ***/

func (handler *CommandHandler) register(
	customers Customers,
	register *commands.Register,
) error {

	recordedEvents := domain.RegisterCustomer(register)

	if err := customers.Register(register.ID(), recordedEvents); err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) confirmEmailAddress(
	customers Customers,
	confirmEmailAddress *commands.ConfirmEmailAddress,
) error {

	customer, err := customers.Of(confirmEmailAddress.ID())
	if err != nil {
		return err
	}

	recordedEvents := customer.ConfirmEmailAddress(confirmEmailAddress)

	if err := customers.Persist(confirmEmailAddress.ID(), recordedEvents); err != nil {
		return err
	}

	for _, event := range recordedEvents {
		switch actualEvent := event.(type) {
		case *events.EmailAddressConfirmationFailed:
			return errors.Mark(errors.New(actualEvent.EventName()), shared.ErrDomainConstraintsViolation)
		}
	}

	return nil
}

func (handler *CommandHandler) changeEmailAddress(
	customers Customers,
	changeEmailAddress *commands.ChangeEmailAddress,
) error {

	customer, err := customers.Of(changeEmailAddress.ID())
	if err != nil {
		return err
	}

	recordedEvents := customer.ChangeEmailAddress(changeEmailAddress)

	if err := customers.Persist(changeEmailAddress.ID(), recordedEvents); err != nil {
		return err
	}

	return nil
}

/*** Command Assertions ***/

func (handler *CommandHandler) assertIsValid(command shared.Command) error {
	if command == nil {
		return errors.New("command is nil interface")
	}

	if reflect.ValueOf(command).IsNil() {
		return errors.Newf("[%s]: command value is nil pointer", command.CommandName())
	}

	if reflect.ValueOf(command.AggregateID()).IsNil() {
		return errors.Newf("[%s]: command was not properly created", command.CommandName())
	}

	return nil
}

func (handler *CommandHandler) assertIsKnown(command shared.Command) error {
	switch command.(type) {
	case *commands.Register, *commands.ConfirmEmailAddress, *commands.ChangeEmailAddress:
		return nil
	default:
		return errors.Newf("[%s] command is unknown", command.CommandName())
	}
}
