package eventsourced

import (
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib"
	"math"

	"github.com/cockroachdb/errors"
)

type Customers struct {
	eventStoreSession lib.EventStore
}

func (customers *Customers) Register(id values.CustomerID, recordedEvents lib.DomainEvents) error {
	streamID := lib.NewStreamID(streamPrefix + "-" + id.ID())

	if err := customers.eventStoreSession.AppendEventsToStream(streamID, recordedEvents); err != nil {
		if errors.Is(err, lib.ErrConcurrencyConflict) {
			return lib.MarkAndWrapError(err, lib.ErrDuplicate, "customers.Register")
		}

		return errors.Wrap(err, "customers.Register")
	}

	return nil
}

func (customers *Customers) EventStream(id values.CustomerID) (lib.DomainEvents, error) {
	streamID := lib.NewStreamID(streamPrefix + "-" + id.ID())

	eventStream, err := customers.eventStoreSession.LoadEventStream(streamID, 0, math.MaxUint32)
	if err != nil {
		return nil, errors.Wrap(err, "customers.EventStream")
	}

	if len(eventStream) == 0 {
		err := errors.New("found empty event stream")
		return nil, lib.MarkAndWrapError(err, lib.ErrNotFound, "customers.EventStream")
	}

	return eventStream, nil
}

func (customers *Customers) Persist(id values.CustomerID, recordedEvents lib.DomainEvents) error {
	streamID := lib.NewStreamID(streamPrefix + "-" + id.ID())

	if err := customers.eventStoreSession.AppendEventsToStream(streamID, recordedEvents); err != nil {
		return errors.Wrap(err, "customers.Persist")
	}

	return nil
}
