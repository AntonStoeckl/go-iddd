package postgres

import (
	"database/sql"
	"strings"

	"github.com/AntonStoeckl/go-iddd/src/shared/es"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/cockroachdb/errors"
	"github.com/lib/pq"
)

type UniqueCustomerEmailAddresses struct {
	uniqueEmailAddressesTableName     string
	buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions
}

func NewUniqueCustomerEmailAddresses(
	uniqueEmailAddressesTableName string,
	buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions,
) *UniqueCustomerEmailAddresses {

	return &UniqueCustomerEmailAddresses{
		uniqueEmailAddressesTableName:     uniqueEmailAddressesTableName,
		buildUniqueEmailAddressAssertions: buildUniqueEmailAddressAssertions,
	}
}

func (s *UniqueCustomerEmailAddresses) AssertUniqueEmailAddress(recordedEvents []es.DomainEvent, tx *sql.Tx) error {
	wrapWithMsg := "assertUniqueEmailAddresse"

	assertions := s.buildUniqueEmailAddressAssertions(recordedEvents...)

	for _, assertion := range assertions {
		switch assertion.DesiredAction() {
		case customer.ShouldAddUniqueEmailAddress:
			if err := s.tryToAdd(assertion.EmailAddressToAdd(), assertion.CustomerID(), tx); err != nil {
				return errors.Wrap(err, wrapWithMsg)
			}
		case customer.ShouldReplaceUniqueEmailAddress:
			if err := s.tryToReplace(assertion.EmailAddressToAdd(), assertion.CustomerID(), tx); err != nil {
				return errors.Wrap(err, wrapWithMsg)
			}
		case customer.ShouldRemoveUniqueEmailAddress:
			if err := s.remove(assertion.CustomerID(), tx); err != nil {
				return errors.Wrap(err, wrapWithMsg)
			}
		}
	}

	return nil
}

func (s *UniqueCustomerEmailAddresses) PurgeUniqueEmailAddress(customerID value.CustomerID, tx *sql.Tx) error {
	return s.remove(customerID, tx)
}

func (s *UniqueCustomerEmailAddresses) tryToAdd(
	emailAddress value.UnconfirmedEmailAddress,
	customerID value.CustomerID,
	tx *sql.Tx,
) error {

	queryTemplate := `INSERT INTO %tablename% VALUES ($1, $2)`
	query := strings.Replace(queryTemplate, "%tablename%", s.uniqueEmailAddressesTableName, 1)

	_, err := tx.Exec(
		query,
		emailAddress.String(),
		customerID.String(),
	)

	if err != nil {
		return s.mapUniqueEmailAddressPostgresErrors(err)
	}

	return nil
}

func (s *UniqueCustomerEmailAddresses) tryToReplace(
	emailAddress value.UnconfirmedEmailAddress,
	customerID value.CustomerID,
	tx *sql.Tx,
) error {

	queryTemplate := `UPDATE %tablename% set email_address = $1 where customer_id = $2`
	query := strings.Replace(queryTemplate, "%tablename%", s.uniqueEmailAddressesTableName, 1)

	_, err := tx.Exec(
		query,
		emailAddress.String(),
		customerID.String(),
	)

	if err != nil {
		return s.mapUniqueEmailAddressPostgresErrors(err)
	}

	return nil
}

func (s *UniqueCustomerEmailAddresses) remove(
	customerID value.CustomerID,
	tx *sql.Tx,
) error {

	queryTemplate := `DELETE FROM %tablename% where customer_id = $1`
	query := strings.Replace(queryTemplate, "%tablename%", s.uniqueEmailAddressesTableName, 1)

	_, err := tx.Exec(
		query,
		customerID.String(),
	)

	if err != nil {
		return s.mapUniqueEmailAddressPostgresErrors(err)
	}

	return nil
}

func (s *UniqueCustomerEmailAddresses) mapUniqueEmailAddressPostgresErrors(err error) error {
	// nolint:errorlint // errors.As() suggested, but somehow cockroachdb/errors can't convert this properly
	if actualErr, ok := err.(*pq.Error); ok {
		if actualErr.Code == "23505" {
			return errors.Mark(errors.New("duplicate email address"), shared.ErrDuplicate)
		}
	}

	return errors.Mark(err, shared.ErrTechnical) // some other DB error (Tx closed, wrong table, ...)
}
