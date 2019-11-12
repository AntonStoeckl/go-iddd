package eventsourced

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"math"

	"golang.org/x/xerrors"
)

type Customers struct {
	eventStoreSession shared.EventStore
	customerFactory   func(eventStream shared.DomainEvents) (*domain.Customer, error)
}

func (customers *Customers) Register(id *values.ID, recordedEvents shared.DomainEvents) error {
	streamID := shared.NewStreamID(streamPrefix + "-" + id.String())

	if err := customers.eventStoreSession.AppendEventsToStream(streamID, recordedEvents); err != nil {
		if xerrors.Is(err, shared.ErrConcurrencyConflict) {
			return xerrors.Errorf("eventSourcedRepositorySession.Register: %s: %w", err, shared.ErrDuplicate)
		}

		return xerrors.Errorf("eventSourcedRepositorySession.Register: %w", err)
	}

	return nil
}

func (customers *Customers) Of(id *values.ID) (*domain.Customer, error) {
	streamID := shared.NewStreamID(streamPrefix + "-" + id.String())

	eventStream, err := customers.eventStoreSession.LoadEventStream(streamID, 0, math.MaxUint32)
	if err != nil {
		return nil, xerrors.Errorf("eventSourcedRepositorySession.Of: %w", err)
	}

	if len(eventStream) == 0 {
		return nil, xerrors.Errorf("eventSourcedRepositorySession.Of: event stream is empty: %w", shared.ErrNotFound)
	}

	customer, err := customers.customerFactory(eventStream)
	if err != nil {
		return nil, xerrors.Errorf("eventSourcedRepositorySession.Of: %w", err)
	}

	return customer, nil
}

func (customers *Customers) Persist(id *values.ID, recordedEvents shared.DomainEvents) error {
	streamID := shared.NewStreamID(streamPrefix + "-" + id.String())

	if err := customers.eventStoreSession.AppendEventsToStream(streamID, recordedEvents); err != nil {
		return xerrors.Errorf("eventSourcedRepositorySession.Persist: %w", err)
	}

	return nil
}
