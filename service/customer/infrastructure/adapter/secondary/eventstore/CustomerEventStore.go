package eventstore

import (
	"database/sql"
	"math"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer"

	"github.com/AntonStoeckl/go-iddd/service/customer/application/command"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
)

const streamPrefix = "customer"

type CustomerEventStore struct {
	eventStore           es.EventStore
	uniqueEmailAddresses command.ForAssertingUniqueEmailAddresses
	db                   *sql.DB
}

func NewCustomerEventStore(
	eventStore es.EventStore,
	uniqueEmailAddresses command.ForAssertingUniqueEmailAddresses,
	db *sql.DB,
) *CustomerEventStore {

	return &CustomerEventStore{
		eventStore:           eventStore,
		uniqueEmailAddresses: uniqueEmailAddresses,
		db:                   db,
	}
}

func (store *CustomerEventStore) RetrieveCustomerEventStream(id values.CustomerID) (es.EventStream, error) {
	wrapWithMsg := "customerEventStore.RetrieveCustomerEventStream"

	eventStream, err := store.eventStore.LoadEventStream(store.streamID(id), 0, math.MaxUint32)
	if err != nil {
		return nil, errors.Wrap(err, wrapWithMsg)
	}

	if len(eventStream) == 0 {
		err := errors.New("customer not found")
		return nil, lib.MarkAndWrapError(err, lib.ErrNotFound, wrapWithMsg)
	}

	return eventStream, nil
}

func (store *CustomerEventStore) RegisterCustomer(recordedEvents es.RecordedEvents, id values.CustomerID) error {
	var err error
	wrapWithMsg := "customerEventStore.RegisterCustomer"

	tx, err := store.db.Begin()
	if err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	uniqueEmailAddressAssertions := customer.BuildUniqueEmailAddressAssertionsFrom(recordedEvents)

	if err = store.uniqueEmailAddresses.Assert(uniqueEmailAddressAssertions, tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = store.eventStore.AppendEventsToStream(store.streamID(id), recordedEvents, tx); err != nil {
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

func (store *CustomerEventStore) AppendToCustomerEventStream(recordedEvents es.RecordedEvents, id values.CustomerID) error {
	var err error
	wrapWithMsg := "customerEventStore.AppendToCustomerEventStream"

	tx, err := store.db.Begin()
	if err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	uniqueEmailAddressAssertions := customer.BuildUniqueEmailAddressAssertionsFrom(recordedEvents)

	if err = store.uniqueEmailAddresses.Assert(uniqueEmailAddressAssertions, tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = store.eventStore.AppendEventsToStream(store.streamID(id), recordedEvents, tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = tx.Commit(); err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	return nil
}

func (store *CustomerEventStore) PurgeCustomerEventStream(id values.CustomerID) error {
	var err error
	wrapWithMsg := "customerEventStore.PurgeCustomerEventStream"

	tx, err := store.db.Begin()
	if err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	if err = store.uniqueEmailAddresses.ClearFor(id, tx); err != nil {
		_ = tx.Rollback()

		return errors.Wrap(err, wrapWithMsg)
	}

	if err = tx.Commit(); err != nil {
		return lib.MarkAndWrapError(err, lib.ErrTechnical, wrapWithMsg)
	}

	if err := store.eventStore.PurgeEventStream(store.streamID(id)); err != nil {
		return errors.Wrap(err, wrapWithMsg)
	}

	return nil
}

func (store *CustomerEventStore) streamID(id values.CustomerID) es.StreamID {
	return es.NewStreamID(streamPrefix + "-" + id.String())
}
