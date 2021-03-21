package shared

import (
	"database/sql"
)

func WrapInTx(
	wrappedFn func(tx *sql.Tx) error,
	db *sql.DB,
) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if err := wrappedFn(tx); err != nil {
		_ = tx.Rollback()

		return err
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()

		return err
	}

	return nil
}
