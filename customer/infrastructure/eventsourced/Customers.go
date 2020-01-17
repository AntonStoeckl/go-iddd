package eventsourced

import (
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"math"

	"golang.org/x/xerrors"
)

type Customers struct {
	eventStoreSession shared.EventStore
}

func (customers *Customers) Register(id values.CustomerID, recordedEvents shared.DomainEvents) error {
	streamID := shared.NewStreamID(streamPrefix + "-" + id.ID())

	if err := customers.eventStoreSession.AppendEventsToStream(streamID, recordedEvents); err != nil {
		if xerrors.Is(err, shared.ErrConcurrencyConflict) {
			return xerrors.Errorf("eventSourcedRepositorySession.Register: %s: %w", err, shared.ErrDuplicate)
		}

		return xerrors.Errorf("eventSourcedRepositorySession.Register: %w", err)
	}

	return nil
}

func (customers *Customers) EventStream(id values.CustomerID) (shared.DomainEvents, error) {
	streamID := shared.NewStreamID(streamPrefix + "-" + id.ID())

	eventStream, err := customers.eventStoreSession.LoadEventStream(streamID, 0, math.MaxUint32)
	if err != nil {
		return nil, xerrors.Errorf("eventSourcedRepositorySession.Of: %w", err)
	}

	if len(eventStream) == 0 {
		return nil, xerrors.Errorf("eventSourcedRepositorySession.Of: event stream is empty: %w", shared.ErrNotFound)
	}

	return eventStream, nil
}

func (customers *Customers) Persist(id values.CustomerID, recordedEvents shared.DomainEvents) error {
	streamID := shared.NewStreamID(streamPrefix + "-" + id.ID())

	if err := customers.eventStoreSession.AppendEventsToStream(streamID, recordedEvents); err != nil {
		return xerrors.Errorf("eventSourcedRepositorySession.Persist: %w", err)
	}

	return nil
}
