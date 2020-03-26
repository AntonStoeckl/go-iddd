package postgres

import (
	"database/sql"
	"strings"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"
	"github.com/lib/pq"
)

type UniqueEmailAddressChecker struct {
	tableName string
}

func NewUniqueEmailAddressChecker(tableName string) *UniqueEmailAddressChecker {
	return &UniqueEmailAddressChecker{tableName: tableName}
}

func (checker *UniqueEmailAddressChecker) AddUniqueEmailAddress(
	emailAddress values.EmailAddress,
	tx *sql.Tx,
) error {

	var err error

	wrapWithMsg := "addUniqueEmailAddress"

	queryTemplate := `INSERT INTO %tablename% VALUES ($1)`
	query := strings.Replace(queryTemplate, "%tablename%", checker.tableName, 1)

	_, err = tx.Exec(
		query,
		emailAddress.EmailAddress(),
	)

	if err != nil {
		return errors.Wrap(checker.mapUniqueEmailAddressErrors(err), wrapWithMsg)
	}

	return nil
}

func (checker *UniqueEmailAddressChecker) ReplaceUniqueEmailAddress(
	previousEmailAddress values.EmailAddress,
	newEmailAddress values.EmailAddress,
	tx *sql.Tx,
) error {

	var err error

	wrapWithMsg := "replaceUniqueEmailAddress"

	queryTemplate := `UPDATE %tablename% set email_address = $1 where email_address = $2`
	query := strings.Replace(queryTemplate, "%tablename%", checker.tableName, 1)

	_, err = tx.Exec(
		query,
		newEmailAddress.EmailAddress(),
		previousEmailAddress.EmailAddress(),
	)

	if err != nil {
		return errors.Wrap(checker.mapUniqueEmailAddressErrors(err), wrapWithMsg)
	}

	return nil
}

func (checker *UniqueEmailAddressChecker) RemoveUniqueEmailAddress(
	newEmailAddress values.EmailAddress,
	tx *sql.Tx) error {

	var err error

	wrapWithMsg := "removeUniqueEmailAddress"

	queryTemplate := `DELETE FROM %tablename% where email_address = $1`
	query := strings.Replace(queryTemplate, "%tablename%", checker.tableName, 1)

	_, err = tx.Exec(
		query,
		newEmailAddress.EmailAddress(),
	)

	if err != nil {
		return errors.Wrap(checker.mapUniqueEmailAddressErrors(err), wrapWithMsg)
	}

	return nil
}

func (checker *UniqueEmailAddressChecker) mapUniqueEmailAddressErrors(err error) error {
	defaultErr := errors.Mark(err, lib.ErrTechnical)

	switch actualErr := err.(type) {
	case *pq.Error:
		switch actualErr.Code {
		case "23505":
			return errors.Mark(errors.Newf("duplicate email address"), lib.ErrDuplicate)
		default:
			return defaultErr // some other postgres error (e.g. table does not exist)
		}
	default:
		return defaultErr // some other DB error (e.g. tx already closed, no connection)
	}
}
