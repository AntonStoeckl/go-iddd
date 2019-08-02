package application

import (
	"database/sql"
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/shared"
	"reflect"

	"github.com/cockroachdb/errors"
)

const maxCommandHandlerRetries = 10

type PersistableCustomers interface {
	domain.Customers
	Persist(customer domain.Customer) error
}

type StartsRepositorySessions interface {
	StartSession(tx *sql.Tx) PersistableCustomers
}

type CommandHandler struct {
	repo StartsRepositorySessions
	db   *sql.DB
}

/*** Factory Method ***/

func NewCommandHandler(repo StartsRepositorySessions, db *sql.DB) *CommandHandler {
	return &CommandHandler{
		repo: repo,
		db:   db,
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

	session := handler.repo.StartSession(tx)

	// call next method in chain
	if err := handler.handleCommand(session, command); err != nil {
		if errTx := tx.Rollback(); errTx != nil {
			return errors.Wrap(err, errTx.Error())
		}

		return err
	}

	if errTx := tx.Commit(); errTx != nil {
		return errors.Mark(errTx, shared.ErrTechnical)
	}

	return nil
}

func (handler *CommandHandler) handleCommand(
	customers PersistableCustomers,
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
	customers PersistableCustomers,
	register *commands.Register,
) error {

	newCustomer := domain.NewCustomerWith(register)

	if err := customers.Register(newCustomer); err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) confirmEmailAddress(
	customers PersistableCustomers,
	confirmEmailAddress *commands.ConfirmEmailAddress,
) error {

	customer, err := customers.Of(confirmEmailAddress.ID())
	if err != nil {
		return err
	}

	if err := customer.Execute(confirmEmailAddress); err != nil {
		return err
	}

	if err := customers.Persist(customer); err != nil {
		return err
	}

	return nil
}

func (handler *CommandHandler) changeEmailAddress(
	customers PersistableCustomers,
	changeEmailAddress *commands.ChangeEmailAddress,
) error {

	customer, err := customers.Of(changeEmailAddress.ID())
	if err != nil {
		return err
	}

	if err := customer.Execute(changeEmailAddress); err != nil {
		return err
	}

	if err := customers.Persist(customer); err != nil {
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
