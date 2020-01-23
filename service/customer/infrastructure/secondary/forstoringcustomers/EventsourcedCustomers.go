package forstoringcustomers

import (
	"database/sql"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib"
	"math"

	"github.com/cockroachdb/errors"
)

const streamPrefix = "customer"

type EventsourcedCustomers struct {
	eventStore lib.EventStore
}

func NewEventsourcedCustomers(eventStore lib.EventStore) *EventsourcedCustomers {
	return &EventsourcedCustomers{eventStore: eventStore}
}

func (customers *EventsourcedCustomers) EventStream(id values.CustomerID) (lib.DomainEvents, error) {
	eventStream, err := customers.eventStore.LoadEventStream(customers.streamID(id), 0, math.MaxUint32)
	if err != nil {
		return nil, errors.Wrap(err, "customers.EventStream")
	}

	if len(eventStream) == 0 {
		err := errors.New("found empty event stream")
		return nil, lib.MarkAndWrapError(err, lib.ErrNotFound, "customers.EventStream")
	}

	return eventStream, nil
}

func (customers *EventsourcedCustomers) Register(id values.CustomerID, recordedEvents lib.DomainEvents, tx *sql.Tx) error {
	if err := customers.eventStore.AppendEventsToStream(customers.streamID(id), recordedEvents, tx); err != nil {
		if errors.Is(err, lib.ErrConcurrencyConflict) {
			err = errors.New("found duplicate customer")
			return lib.MarkAndWrapError(err, lib.ErrDuplicate, "customers.Register")
		}

		return errors.Wrap(err, "customers.Register")
	}

	return nil
}

func (customers *EventsourcedCustomers) Persist(id values.CustomerID, recordedEvents lib.DomainEvents, tx *sql.Tx) error {
	if err := customers.eventStore.AppendEventsToStream(customers.streamID(id), recordedEvents, tx); err != nil {
		return errors.Wrap(err, "customers.Persist")
	}

	return nil
}

func (customers *EventsourcedCustomers) Delete(id values.CustomerID) error {
	if err := customers.eventStore.PurgeEventStream(customers.streamID(id)); err != nil {
		return errors.Wrap(err, "customers.Delete")
	}

	return nil
}

func (customers *EventsourcedCustomers) streamID(id values.CustomerID) lib.StreamID {
	return lib.NewStreamID(streamPrefix + "-" + id.ID())
}
