package postgres

import (
	"database/sql"
	"math"
	"strings"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
	"github.com/lib/pq"
)

const streamPrefix = "customer"

type CustomerEventStore struct {
	db                            *sql.DB
	eventStoreTableName           string
	uniqueEmailAddressesTableName string
	marshalDomainEvent            es.MarshalDomainEvent
	unmarshalDomainEvent          es.UnmarshalDomainEvent
}

func NewCustomerEventStore(
	db *sql.DB,
	eventStoreTableName string,
	uniqueEmailAddressesTableName string,
	marshalDomainEvent es.MarshalDomainEvent,
	unmarshalDomainEvent es.UnmarshalDomainEvent,
) *CustomerEventStore {

	return &CustomerEventStore{
		db:                            db,
		eventStoreTableName:           eventStoreTableName,
		uniqueEmailAddressesTableName: uniqueEmailAddressesTableName,
		marshalDomainEvent:            marshalDomainEvent,
		unmarshalDomainEvent:          unmarshalDomainEvent,
	}
}

func (s *CustomerEventStore) RetrieveCustomerEventStream(id values.CustomerID) (es.EventStream, error) {
	wrapWithMsg := "customerEventStore.RetrieveCustomerEventStream"

	eventStream, err := s.loadEventStream(s.streamID(id), 0, math.MaxUint32)
	if err != nil {
		return nil, errors.Wrap(err, wrapWithMsg)
	}

	if len(eventStream) == 0 {
		err := errors.New("customer not found")
		return nil, lib.MarkAndWrapError(err, lib.ErrNotFound, wrapWithMsg)
	}

	return eventStream, nil
}

func (s *CustomerEventStore) RegisterCustomer(recordedEvents es.RecordedEvents, id values.CustomerID) error {
	var err error
	wrapWithMsg := "customerEventStore.RegisterCustomer"

	tx, err := s.db.Begin()
	if err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	uniqueEmailAddressAssertions := customer.BuildUniqueEmailAddressAssertionsFrom(recordedEvents)

	if err = s.assertUniqueEmailAddress(uniqueEmailAddressAssertions, tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = s.appendEventsToStream(s.streamID(id), recordedEvents, tx); err != nil {
		_ = tx.Rollback()

		if errors.Is(err, lib.ErrConcurrencyConflict) {
			return lib.MarkAndWrapError(errors.New("found duplicate customer"), lib.ErrDuplicate, wrapWithMsg)
		}

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = tx.Commit(); err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	return nil
}

func (s *CustomerEventStore) AppendToCustomerEventStream(recordedEvents es.RecordedEvents, id values.CustomerID) error {
	var err error
	wrapWithMsg := "customerEventStore.AppendToCustomerEventStream"

	tx, err := s.db.Begin()
	if err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	uniqueEmailAddressAssertions := customer.BuildUniqueEmailAddressAssertionsFrom(recordedEvents)

	if err = s.assertUniqueEmailAddress(uniqueEmailAddressAssertions, tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = s.appendEventsToStream(s.streamID(id), recordedEvents, tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = tx.Commit(); err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	return nil
}

func (s *CustomerEventStore) PurgeCustomerEventStream(id values.CustomerID) error {
	var err error
	wrapWithMsg := "customerEventStore.PurgeCustomerEventStream"

	tx, err := s.db.Begin()
	if err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	if err = s.clearUniqueEmailAddress(id, tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = tx.Commit(); err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	if err := s.purgeEventStream(s.streamID(id)); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}

func (s *CustomerEventStore) streamID(id values.CustomerID) es.StreamID {
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
		return nil, lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	var eventStream es.EventStream
	var eventName string
	var payload string
	var streamVersion uint
	var domainEvent es.DomainEvent

	for eventRows.Next() {
		if eventRows.Err() != nil {
			return nil, lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
		}

		if err = eventRows.Scan(&eventName, &payload, &streamVersion); err != nil {
			return nil, lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
		}

		if domainEvent, err = s.unmarshalDomainEvent(eventName, []byte(payload), streamVersion); err != nil {
			return nil, lib.MarkAndWrapError(err, lib.ErrUnmarshalingFailed, wrapWithMsg)
		}

		eventStream = append(eventStream, domainEvent)
	}

	return eventStream, nil
}

func (s *CustomerEventStore) appendEventsToStream(
	streamID es.StreamID,
	events es.RecordedEvents,
	tx *sql.Tx,
) error {

	var err error
	wrapWithMsg := "appendEventsToStream"

	queryTemplate := `INSERT INTO %name% (stream_id, stream_version, event_name, occurred_at, payload)
						VALUES ($1, $2, $3, $4, $5)`
	query := strings.Replace(queryTemplate, "%name%", s.eventStoreTableName, 1)

	for _, event := range events {
		var eventJson []byte

		eventJson, err = s.marshalDomainEvent(event)
		if err != nil {
			return lib.MarkAndWrapError(err, lib.ErrMarshalingFailed, wrapWithMsg)
		}

		_, err = tx.Exec(
			query,
			streamID.String(),
			event.Meta().StreamVersion(),
			event.Meta().EventName(),
			event.Meta().OccurredAt(),
			eventJson,
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
		return lib.MarkAndWrapError(err, lib.ErrTechnical, "purgeEventStream")
	}

	return nil
}

func (s *CustomerEventStore) mapEventStorePostgresErrors(err error) error {
	switch actualErr := err.(type) {
	case *pq.Error:
		switch actualErr.Code {
		case "23505":
			return errors.Mark(err, lib.ErrConcurrencyConflict)
		}
	}

	return errors.Mark(err, lib.ErrTechnical) // some other DB error (Tx closed, wrong table, ...)
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

func (s *CustomerEventStore) clearUniqueEmailAddress(customerID values.CustomerID, tx *sql.Tx) error {
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
	emailAddress values.EmailAddress,
	customerID values.CustomerID,
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
	previousEmailAddress values.EmailAddress,
	newEmailAddress values.EmailAddress,
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
	newEmailAddress values.EmailAddress,
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
	switch actualErr := err.(type) {
	case *pq.Error:
		switch actualErr.Code {
		case "23505":
			return errors.Mark(errors.Newf("duplicate email address"), lib.ErrDuplicate)
		}
	}

	return errors.Mark(err, lib.ErrTechnical) // some other DB error (Tx closed, wrong table, ...)
}
