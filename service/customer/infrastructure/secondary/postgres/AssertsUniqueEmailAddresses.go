package postgres

import (
	"database/sql"
	"strings"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
	"github.com/lib/pq"
)

type AssertsUniqueEmailAddresses struct {
	tableName string
}

func NewAssertsUniqueEmailAddresses(tableName string) *AssertsUniqueEmailAddresses {
	return &AssertsUniqueEmailAddresses{tableName: tableName}
}

func (asserter *AssertsUniqueEmailAddresses) Assert(recordedEvents es.DomainEvents, tx *sql.Tx) error {
	wrapWithMsg := "assertsUniqueEmailAddresses"

	for _, event := range recordedEvents {
		switch actualEvent := event.(type) {
		case events.CustomerRegistered:
			if err := asserter.add(actualEvent.EmailAddress(), actualEvent.CustomerID(), tx); err != nil {
				return errors.Wrap(err, wrapWithMsg)
			}
		case events.CustomerEmailAddressChanged:
			if err := asserter.replace(actualEvent.PreviousEmailAddress(), actualEvent.EmailAddress(), tx); err != nil {
				return errors.Wrap(err, wrapWithMsg)
			}
		case events.CustomerDeleted:
			if err := asserter.remove(actualEvent.EmailAddress(), tx); err != nil {
				return errors.Wrap(err, wrapWithMsg)
			}
		}
	}

	return nil
}

func (asserter *AssertsUniqueEmailAddresses) Remove(customerID values.CustomerID, tx *sql.Tx) error {
	queryTemplate := `DELETE FROM %tablename% WHERE customer_id = $1`
	query := strings.Replace(queryTemplate, "%tablename%", asserter.tableName, 1)

	_, err := tx.Exec(
		query,
		customerID.ID(),
	)

	if err != nil {
		return asserter.mapUniqueEmailAddressErrors(err)
	}

	return nil
}

func (asserter *AssertsUniqueEmailAddresses) add(
	emailAddress values.EmailAddress,
	customerID values.CustomerID,
	tx *sql.Tx,
) error {

	queryTemplate := `INSERT INTO %tablename% VALUES ($1, $2)`
	query := strings.Replace(queryTemplate, "%tablename%", asserter.tableName, 1)

	_, err := tx.Exec(
		query,
		emailAddress.EmailAddress(),
		customerID.ID(),
	)

	if err != nil {
		return asserter.mapUniqueEmailAddressErrors(err)
	}

	return nil
}

func (asserter *AssertsUniqueEmailAddresses) replace(
	previousEmailAddress values.EmailAddress,
	newEmailAddress values.EmailAddress,
	tx *sql.Tx,
) error {

	queryTemplate := `UPDATE %tablename% set email_address = $1 where email_address = $2`
	query := strings.Replace(queryTemplate, "%tablename%", asserter.tableName, 1)

	_, err := tx.Exec(
		query,
		newEmailAddress.EmailAddress(),
		previousEmailAddress.EmailAddress(),
	)

	if err != nil {
		return asserter.mapUniqueEmailAddressErrors(err)
	}

	return nil
}

func (asserter *AssertsUniqueEmailAddresses) remove(
	newEmailAddress values.EmailAddress,
	tx *sql.Tx,
) error {

	queryTemplate := `DELETE FROM %tablename% where email_address = $1`
	query := strings.Replace(queryTemplate, "%tablename%", asserter.tableName, 1)

	_, err := tx.Exec(
		query,
		newEmailAddress.EmailAddress(),
	)

	if err != nil {
		return asserter.mapUniqueEmailAddressErrors(err)
	}

	return nil
}

func (asserter *AssertsUniqueEmailAddresses) mapUniqueEmailAddressErrors(err error) error {
	switch actualErr := err.(type) {
	case *pq.Error:
		switch actualErr.Code {
		case "23505":
			return errors.Mark(errors.Newf("duplicate email address"), lib.ErrDuplicate)
		}
	}

	return errors.Mark(err, lib.ErrTechnical) // some other DB error (Tx closed, wrong table, ...)
}
