package postgres

import (
	"database/sql"
	"strings"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/identity/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
	"github.com/lib/pq"
)

type UniqueIdentities struct {
	tableName string
	db        *sql.DB
	tx        *sql.Tx
}

func NewUniqueIdentities(tableName string, db *sql.DB) *UniqueIdentities {
	return &UniqueIdentities{
		tableName: tableName,
		db:        db,
	}
}

func (i *UniqueIdentities) WithTx(tx *sql.Tx) application.ForStoringUniqueIdentities {
	return &UniqueIdentities{
		tableName: i.tableName,
		db:        i.db,
		tx:        tx,
	}
}

func (i *UniqueIdentities) FindIdentity(emailAddress value.UnconfirmedEmailAddress) (value.IdentityID, error) {
	wrapWithMsg := "findIdentity"

	queryTemplate := `SELECT identity_id FROM %tablename% WHERE email_address = $1`
	query := strings.Replace(queryTemplate, "%tablename%", i.tableName, 1)

	row := i.db.QueryRow(query, emailAddress.String())

	var identityID string

	err := row.Scan(&identityID)
	if row.Err() != nil {
		return "", shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	return value.RebuildIdentityID(identityID), nil
}

func (i *UniqueIdentities) AddIdentity(identityID value.IdentityID, emailAddress value.UnconfirmedEmailAddress) error {
	queryTemplate := `INSERT INTO %tablename% VALUES ($1, $2)`
	query := strings.Replace(queryTemplate, "%tablename%", i.tableName, 1)

	_, err := i.tx.Exec(
		query,
		emailAddress.String(),
		identityID.String(),
	)

	if err != nil {
		return i.mapToPostgresError(err)
	}

	return nil
}

func (i *UniqueIdentities) RemoveIdentity(identityID value.IdentityID) error {
	queryTemplate := `DELETE FROM %tablename% WHERE identity_id = $1`
	query := strings.Replace(queryTemplate, "%tablename%", i.tableName, 1)

	_, err := i.db.Exec(
		query,
		identityID.String(),
	)

	if err != nil {
		return i.mapToPostgresError(err)
	}

	return nil
}

func (i *UniqueIdentities) mapToPostgresError(err error) error {
	// nolint:errorlint // errors.As() suggested, but somehow cockroachdb/errors can't convert this properly
	if actualErr, ok := err.(*pq.Error); ok {
		if actualErr.Code == "23505" {
			return errors.Mark(errors.New("duplicate email address"), shared.ErrDuplicate)
		}
	}

	return errors.Mark(err, shared.ErrTechnical) // some other DB error (Tx closed, wrong table, ...)
}
