package application

import (
	"database/sql"
	"go-iddd/shared"
)

type TxCommandHandler struct {
	db           *sql.DB
	innerHandler CommandHandler
}

/*** Factory Method ***/

func NewTxCommandHandler(db *sql.DB, innerHandler CommandHandler) *TxCommandHandler {
	return &TxCommandHandler{
		db:           db,
		innerHandler: innerHandler,
	}
}

/*** Implement shared.CommandHandler ***/

func (handler *TxCommandHandler) Handle(command shared.Command) error {
	tx, err := handler.db.Begin()
	if err != nil {
		return err // TODO: map error
	}

	if err := handler.innerHandler.Handle(command); err != nil {
		if err := tx.Rollback(); err != nil {
			return err // TODO: map error
		}

		return err // TODO: map error
	}

	if err := tx.Commit(); err != nil {
		return err // TODO: map error
	}

	return nil
}
