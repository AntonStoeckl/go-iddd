package postgres

import (
	"database/sql"
	"math"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
)

const streamPrefix = "customer"

type forRetrievingEventStreams func(streamID es.StreamID, fromVersion uint, maxEvents uint, db *sql.DB) (es.EventStream, error)
type forAppendingEventsToStreams func(streamID es.StreamID, events []es.DomainEvent, tx *sql.Tx) error
type forPurgingEventStreams func(streamID es.StreamID, tx *sql.Tx) error
type forAssertingUniqueEmailAddresses func(recordedEvents []es.DomainEvent, tx *sql.Tx) error
type forPurgingUniqueEmailAddresses func(customerID value.CustomerID, tx *sql.Tx) error

type CustomerPostgresEventStore struct {
	db                       *sql.DB
	retrieveEventStream      forRetrievingEventStreams
	appendEventsToStream     forAppendingEventsToStreams
	purgeEventStream         forPurgingEventStreams
	assertUniqueEmailAddress forAssertingUniqueEmailAddresses
	purgeUniqueEmailAddress  forPurgingUniqueEmailAddresses
}

func NewCustomerPostgresEventStore(
	db *sql.DB,
	retrieveEventStream forRetrievingEventStreams,
	appendEventsToStream forAppendingEventsToStreams,
	purgeEventStream forPurgingEventStreams,
	assertUniqueEmailAddress forAssertingUniqueEmailAddresses,
	purgeUniqueEmailAddress forPurgingUniqueEmailAddresses,
) application.EventStoreInterface {

	return &CustomerPostgresEventStore{
		db:                       db,
		retrieveEventStream:      retrieveEventStream,
		appendEventsToStream:     appendEventsToStream,
		purgeEventStream:         purgeEventStream,
		assertUniqueEmailAddress: assertUniqueEmailAddress,
		purgeUniqueEmailAddress:  purgeUniqueEmailAddress,
	}
}

func (s *CustomerPostgresEventStore) RetrieveEventStream(id value.CustomerID) (es.EventStream, error) {
	wrapWithMsg := "CustomerPostgresEventStore.RetrieveEventStream"

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

func (s *CustomerPostgresEventStore) StartEventStream(customerRegistered domain.CustomerRegistered) error {
	var err error
	wrapWithMsg := "CustomerPostgresEventStore.StartEventStream"

	tx, err := s.db.Begin()
	if err != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	recordedEvents := []es.DomainEvent{customerRegistered}

	if err = s.assertUniqueEmailAddress(recordedEvents, tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	streamID := s.streamID(customerRegistered.CustomerID())

	if err = s.appendEventsToStream(streamID, recordedEvents, tx); err != nil {
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

func (s *CustomerPostgresEventStore) AppendToEventStream(recordedEvents es.RecordedEvents, id value.CustomerID) error {
	var err error
	wrapWithMsg := "CustomerPostgresEventStore.AppendToEventStream"

	tx, err := s.db.Begin()
	if err != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	if err = s.assertUniqueEmailAddress(recordedEvents, tx); err != nil {
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

func (s *CustomerPostgresEventStore) PurgeEventStream(id value.CustomerID) error {
	var err error
	wrapWithMsg := "CustomerPostgresEventStore.PurgeEventStream"

	tx, err := s.db.Begin()
	if err != nil {
		return shared.MarkAndWrapError(err, shared.ErrTechnical, wrapWithMsg)
	}

	if err = s.purgeUniqueEmailAddress(id, tx); err != nil {
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

func (s *CustomerPostgresEventStore) streamID(id value.CustomerID) es.StreamID {
	return es.BuildStreamID(streamPrefix + "-" + id.String())
}
