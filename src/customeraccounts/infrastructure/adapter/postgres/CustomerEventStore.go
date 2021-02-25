package postgres

import (
	"database/sql"
	"math"
	"strings"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	"github.com/lib/pq"
)

const streamPrefix = "customer"

type forRetrievingEventStreams func(streamID es.StreamID, fromVersion uint, maxEvents uint, db *sql.DB) (es.EventStream, error)
type forAppendingEventsToStreams func(streamID es.StreamID, events []es.DomainEvent, tx *sql.Tx) error
type forPurgingEventStreams func(streamID es.StreamID, tx *sql.Tx) error

type CustomerEventStore struct {
	db                                *sql.DB
	retrieveEventStream               forRetrievingEventStreams
	appendEventsToStream              forAppendingEventsToStreams
	purgeEventStream                  forPurgingEventStreams
	uniqueEmailAddressesTableName     string
	buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions
}

func NewCustomerEventStore(
	db *sql.DB,
	eventStore *EventStore,
	uniqueEmailAddressesTableName string,
	buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions,
) *CustomerEventStore {

	return &CustomerEventStore{
		db:                                db,
		retrieveEventStream:               eventStore.RetrieveEventStream,
		appendEventsToStream:              eventStore.AppendEventsToStream,
		purgeEventStream:                  eventStore.PurgeEventStream,
		uniqueEmailAddressesTableName:     uniqueEmailAddressesTableName,
		buildUniqueEmailAddressAssertions: buildUniqueEmailAddressAssertions,
	}
}

func (s *CustomerEventStore) RetrieveEventStream(id value.CustomerID) (es.EventStream, error) {
	wrapWithMsg := "customerEventStore.RetrieveEventStream"

	eventStream, err := s.retrieveEventStream(s.streamID(id), 0, math.MaxUint32, s.db)
	if err != nil {
		return nil, errors.Wrap(err, wrapWithMsg)
	}

	if len(eventStream) == 0 {
		err := errors.New("customer not found")
		return nil, shared.MarkAndWrapError(err, shared.ErrNotFound, wrapWithMsg)
	}

	return eventStream, nil
}

func (s *CustomerEventStore) StartEventStream(customerRegistered domain.CustomerRegistered) error {
	var err error
	wrapWithMsg := "customerEventStore.StartEventStream"

	tx, err := s.db.Begin()
	if err != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	assertionsForUniqueEmailAddresses := s.buildUniqueEmailAddressAssertions(customerRegistered)

	if err = s.assertUniqueEmailAddress(assertionsForUniqueEmailAddresses, tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	events := []es.DomainEvent{customerRegistered}
	streamID := s.streamID(customerRegistered.CustomerID())

	if err = s.appendEventsToStream(streamID, events, tx); err != nil {
		_ = tx.Rollback()

		if errors.Is(err, shared.ErrConcurrencyConflict) {
			return shared.MarkAndWrapError(errors.New("found duplicate customer"), shared.ErrDuplicate, wrapWithMsg)
		}

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = tx.Commit(); err != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	return nil
}

func (s *CustomerEventStore) AppendToEventStream(recordedEvents es.RecordedEvents, id value.CustomerID) error {
	var err error
	wrapWithMsg := "customerEventStore.AppendToEventStream"

	tx, err := s.db.Begin()
	if err != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	assertionsForUniqueEmailAddresses := s.buildUniqueEmailAddressAssertions(recordedEvents...)

	if err = s.assertUniqueEmailAddress(assertionsForUniqueEmailAddresses, tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = s.appendEventsToStream(s.streamID(id), recordedEvents, tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = tx.Commit(); err != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	return nil
}

func (s *CustomerEventStore) PurgeEventStream(id value.CustomerID) error {
	var err error
	wrapWithMsg := "customerEventStore.PurgeEventStream"

	tx, err := s.db.Begin()
	if err != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	if err = s.remove(id, tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = s.purgeEventStream(s.streamID(id), tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = tx.Commit(); err != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	return nil
}

func (s *CustomerEventStore) streamID(id value.CustomerID) es.StreamID {
	return es.BuildStreamID(streamPrefix + "-" + id.String())
}

/***** local methods for asserting unique email addresses *****/

func (s *CustomerEventStore) assertUniqueEmailAddress(assertions customer.UniqueEmailAddressAssertions, tx *sql.Tx) error {
	wrapWithMsg := "assertUniqueEmailAddresse"

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

func (s *CustomerEventStore) tryToAdd(
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

func (s *CustomerEventStore) tryToReplace(
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

func (s *CustomerEventStore) remove(
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

func (s *CustomerEventStore) mapUniqueEmailAddressPostgresErrors(err error) error {
	// nolint:errorlint // errors.As() suggested, but somehow cockroachdb/errors can't convert this properly
	if actualErr, ok := err.(*pq.Error); ok {
		if actualErr.Code == "23505" {
			return errors.Mark(errors.New("duplicate email address"), shared.ErrDuplicate)
		}
	}

	return errors.Mark(err, shared.ErrTechnical) // some other DB error (Tx closed, wrong table, ...)
}
