package customers

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"math"

	"golang.org/x/xerrors"
)

type EventSourcedRepositorySession struct {
	eventStoreSession shared.EventStore
	customerFactory   func(eventStream shared.DomainEvents) (domain.Customer, error)
	identityMap       *IdentityMap
}

/***** Implement domain.Customers *****/

func (session *EventSourcedRepositorySession) Register(customer domain.Customer) error {
	streamID := shared.NewStreamID(streamPrefix + "-" + customer.ID().String())

	if err := session.eventStoreSession.AppendEventsToStream(streamID, customer.RecordedEvents(true)); err != nil {
		if xerrors.Is(err, shared.ErrConcurrencyConflict) {
			return xerrors.Errorf("eventSourcedRepositorySession.Register: %s: %w", err, shared.ErrDuplicate)
		}

		return xerrors.Errorf("eventSourcedRepositorySession.Register: %w", err)
	}

	return nil
}

func (session *EventSourcedRepositorySession) Of(id *values.ID) (domain.Customer, error) {
	streamID := shared.NewStreamID(streamPrefix + "-" + id.String())

	if memoizedCustomer, found := session.identityMap.MemoizedCustomerOf(id); found {
		latestEvents, err := session.eventStoreSession.LoadEventStream(
			streamID,
			memoizedCustomer.StreamVersion(),
			math.MaxUint32,
		)
		if err != nil {
			return nil, xerrors.Errorf("eventSourcedRepositorySession.Of: %w", err)
		}

		memoizedCustomer.Apply(latestEvents)

		session.identityMap.Memoize(memoizedCustomer)

		return memoizedCustomer, nil
	}

	eventStream, err := session.eventStoreSession.LoadEventStream(streamID, 0, math.MaxUint32)
	if err != nil {
		return nil, xerrors.Errorf("eventSourcedRepositorySession.Of: %w", err)
	}

	if len(eventStream) == 0 {
		return nil, xerrors.Errorf("eventSourcedRepositorySession.Of: event stream is empty: %w", shared.ErrNotFound)
	}

	customer, err := session.customerFactory(eventStream)
	if err != nil {
		return nil, xerrors.Errorf("eventSourcedRepositorySession.Of: %w", err)
	}

	session.identityMap.Memoize(customer)

	return customer, nil
}

/***** Implement application.PersistsCustomers *****/

func (session *EventSourcedRepositorySession) Persist(customer domain.Customer) error {
	streamID := shared.NewStreamID(streamPrefix + "-" + customer.ID().String())

	if err := session.eventStoreSession.AppendEventsToStream(streamID, customer.RecordedEvents(true)); err != nil {
		return xerrors.Errorf("eventSourcedRepositorySession.Persist: %w", err)
	}

	return nil
}
