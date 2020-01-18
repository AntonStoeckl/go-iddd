package eventsourced

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"math"

	"github.com/cockroachdb/errors"
)

type Customers struct {
	eventStoreSession shared.EventStore
}

func (customers *Customers) Register(id values.CustomerID, recordedEvents shared.DomainEvents) error {
	streamID := shared.NewStreamID(streamPrefix + "-" + id.ID())

	if err := customers.eventStoreSession.AppendEventsToStream(streamID, recordedEvents); err != nil {
		if errors.Is(err, shared.ErrConcurrencyConflict) {
			return shared.MarkAndWrapError(err, shared.ErrDuplicate, "customers.Register")
		}

		return errors.Wrap(err, "customers.Register")
	}

	return nil
}

func (customers *Customers) EventStream(id values.CustomerID) (shared.DomainEvents, error) {
	streamID := shared.NewStreamID(streamPrefix + "-" + id.ID())

	eventStream, err := customers.eventStoreSession.LoadEventStream(streamID, 0, math.MaxUint32)
	if err != nil {
		return nil, errors.Wrap(err, "customers.EventStream")
	}

	if len(eventStream) == 0 {
		err := errors.New("found empty event stream")
		return nil, shared.MarkAndWrapError(err, shared.ErrNotFound, "customers.EventStream")
	}

	return eventStream, nil
}

func (customers *Customers) Persist(id values.CustomerID, recordedEvents shared.DomainEvents) error {
	streamID := shared.NewStreamID(streamPrefix + "-" + id.ID())

	if err := customers.eventStoreSession.AppendEventsToStream(streamID, recordedEvents); err != nil {
		return errors.Wrap(err, "customers.Persist")
	}

	return nil
}
