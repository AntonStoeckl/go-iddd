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

type CustomerEventStore struct {
	db                                *sql.DB
	eventStoreTableName               string
	marshalDomainEvent                es.MarshalDomainEvent
	unmarshalDomainEvent              es.UnmarshalDomainEvent
	uniqueEmailAddressesTableName     string
	buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions
}

func NewCustomerEventStore(
	db *sql.DB,
	eventStoreTableName string,
	marshalDomainEvent es.MarshalDomainEvent,
	unmarshalDomainEvent es.UnmarshalDomainEvent,
	uniqueEmailAddressesTableName string,
	buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions,
) *CustomerEventStore {

	return &CustomerEventStore{
		db:                                db,
		eventStoreTableName:               eventStoreTableName,
		marshalDomainEvent:                marshalDomainEvent,
		unmarshalDomainEvent:              unmarshalDomainEvent,
		uniqueEmailAddressesTableName:     uniqueEmailAddressesTableName,
		buildUniqueEmailAddressAssertions: buildUniqueEmailAddressAssertions,
	}
}

func (s *CustomerEventStore) RetrieveEventStream(id value.CustomerID) (es.EventStream, error) {
	wrapWithMsg := "customerEventStore.RetrieveEventStream"

	eventStream, err := s.loadEventStream(s.streamID(id), 0, math.MaxUint32)
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

	if err = s.appendEventsToStream(tx, s.streamID(customerRegistered.CustomerID()), customerRegistered); err != nil {
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

	if err = s.appendEventsToStream(tx, s.streamID(id), recordedEvents...); err != nil {
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

	if err = s.clearUniqueEmailAddress(id, tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = tx.Commit(); err != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	if err := s.purgeEventStream(s.streamID(id)); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}

func (s *CustomerEventStore) streamID(id value.CustomerID) es.StreamID {
	return es.NewStreamID(streamPrefix + "-" + id.String())
}

/***** local methods for reading from and writing to the event store *****/

func (s *CustomerEventStore) loadEventStream(
	streamID es.StreamID,
	fromVersion uint,
	maxEvents uint,
) (es.EventStream, error) {

	var err error
	wrapWithMsg := "loadEventStream"

	queryTemplate := `SELECT event_name, payload, stream_version FROM %name% 
						WHERE stream_id = $1 AND stream_version >= $2
						ORDER BY stream_version ASC
						LIMIT $3`

	query := strings.Replace(queryTemplate, "%name%", s.eventStoreTableName, 1)

	eventRows, err := s.db.Query(query, streamID.String(), fromVersion, maxEvents)
	if err != nil {
		return nil, shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	defer eventRows.Close()

	var eventStream es.EventStream
	var eventName string
	var payload string
	var streamVersion uint
	var domainEvent es.DomainEvent

	for eventRows.Next() {
		if eventRows.Err() != nil {
			return nil, shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
		}

		if err = eventRows.Scan(&eventName, &payload, &streamVersion); err != nil {
			return nil, shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
		}

		if domainEvent, err = s.unmarshalDomainEvent(eventName, []byte(payload), streamVersion); err != nil {
			return nil, shared.MarkAndWrapError(err, shared.ErrUnmarshalingFailed, wrapWithMsg)
		}

		eventStream = append(eventStream, domainEvent)
	}

	return eventStream, nil
}

func (s *CustomerEventStore) appendEventsToStream(
	tx *sql.Tx,
	streamID es.StreamID,
	events ...es.DomainEvent,
) error {

	var err error
	wrapWithMsg := "appendEventsToStream"

	queryTemplate := `INSERT INTO %name% (stream_id, stream_version, event_name, occurred_at, payload)
						VALUES ($1, $2, $3, $4, $5)`
	query := strings.Replace(queryTemplate, "%name%", s.eventStoreTableName, 1)

	for _, event := range events {
		var eventJSON []byte

		eventJSON, err = s.marshalDomainEvent(event)
		if err != nil {
			return shared.MarkAndWrapError(err, shared.ErrMarshalingFailed, wrapWithMsg)
		}

		_, err = tx.Exec(
			query,
			streamID.String(),
			event.Meta().StreamVersion(),
			event.Meta().EventName(),
			event.Meta().OccurredAt(),
			eventJSON,
		)

		if err != nil {
			return errors.Wrap(s.mapEventStorePostgresErrors(err), wrapWithMsg)
		}
	}

	return nil
}

func (s *CustomerEventStore) purgeEventStream(streamID es.StreamID) error {
	queryTemplate := `DELETE FROM %name% WHERE stream_id = $1`
	query := strings.Replace(queryTemplate, "%name%", s.eventStoreTableName, 1)

	if _, err := s.db.Exec(query, streamID.String()); err != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, "purgeEventStream")
	}

	return nil
}

func (s *CustomerEventStore) mapEventStorePostgresErrors(err error) error {
	// nolint:errorlint // errors.As() suggested, but somehow cockroachdb/errors can't convert this properly
	if actualErr, ok := err.(*pq.Error); ok {
		if actualErr.Code == "23505" {
			return errors.Mark(err, shared.ErrConcurrencyConflict)
		}
	}

	return errors.Mark(err, shared.ErrTechnical) // some other DB error (Tx closed, wrong table, ...)
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
			if err := s.tryToReplace(assertion.EmailAddressToRemove(), assertion.EmailAddressToAdd(), tx); err != nil {
				return errors.Wrap(err, wrapWithMsg)
			}
		case customer.ShouldRemoveUniqueEmailAddress:
			if err := s.remove(assertion.EmailAddressToRemove(), tx); err != nil {
				return errors.Wrap(err, wrapWithMsg)
			}
		}
	}

	return nil
}

func (s *CustomerEventStore) clearUniqueEmailAddress(customerID value.CustomerID, tx *sql.Tx) error {
	queryTemplate := `DELETE FROM %tablename% WHERE customer_id = $1`
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

func (s *CustomerEventStore) tryToAdd(
	emailAddress value.EmailAddress,
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
	previousEmailAddress value.EmailAddress,
	newEmailAddress value.EmailAddress,
	tx *sql.Tx,
) error {

	queryTemplate := `UPDATE %tablename% set email_address = $1 where email_address = $2`
	query := strings.Replace(queryTemplate, "%tablename%", s.uniqueEmailAddressesTableName, 1)

	_, err := tx.Exec(
		query,
		newEmailAddress.String(),
		previousEmailAddress.String(),
	)

	if err != nil {
		return s.mapUniqueEmailAddressPostgresErrors(err)
	}

	return nil
}

func (s *CustomerEventStore) remove(
	newEmailAddress value.EmailAddress,
	tx *sql.Tx,
) error {

	queryTemplate := `DELETE FROM %tablename% where email_address = $1`
	query := strings.Replace(queryTemplate, "%tablename%", s.uniqueEmailAddressesTableName, 1)

	_, err := tx.Exec(
		query,
		newEmailAddress.String(),
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
			return errors.Mark(errors.Newf("duplicate email address"), shared.ErrDuplicate)
		}
	}

	return errors.Mark(err, shared.ErrTechnical) // some other DB error (Tx closed, wrong table, ...)
}
