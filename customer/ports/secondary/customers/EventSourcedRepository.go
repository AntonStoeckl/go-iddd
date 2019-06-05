package customers

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"

	"golang.org/x/xerrors"
)

type customerFactory func(eventStream shared.DomainEvents) (domain.Customer, error)

type eventSourcedRepository struct {
	eventStore      shared.EventStore
	customerFactory customerFactory
}

func NewEventSourcedRepository(
	eventStore shared.EventStore,
	customerFactory customerFactory,
) *eventSourcedRepository {

	return &eventSourcedRepository{
		eventStore:      eventStore,
		customerFactory: customerFactory,
	}
}

/***** Implement domain.Customers *****/

func (repo *eventSourcedRepository) Register(customer domain.Customer) error {
	if err := repo.eventStore.AppendToStream(customer.AggregateID(), customer.RecordedEvents()); err != nil {
		if xerrors.Is(err, shared.ErrConcurrencyConflict) {
			return xerrors.Errorf("customers[eventSourcedRepository].Register: %s: %w", err, shared.ErrDuplicate)
		}

		return xerrors.Errorf("customers[eventSourcedRepository].Register: %w", err)
	}

	return nil
}

func (repo *eventSourcedRepository) Of(id *values.ID) (domain.Customer, error) {
	eventStream, err := repo.eventStore.LoadEventStream(id)
	if err != nil {
		return nil, xerrors.Errorf("customers[eventSourcedRepository].Of: %w", err)
	}

	if len(eventStream) == 0 {
		return nil, xerrors.Errorf("customers[eventSourcedRepository].Of: event stream is empty: %w", shared.ErrNotFound)
	}

	customer, err := repo.customerFactory(eventStream)
	if err != nil {
		return nil, xerrors.Errorf("customers[eventSourcedRepository].Of: %w", err)
	}

	return customer, nil
}

/***** Implement shared.PersistsEventRecordingAggregates *****/

func (repo *eventSourcedRepository) Persist(aggregate shared.EventRecordingAggregate) error {
	err := repo.eventStore.AppendToStream(
		aggregate.AggregateID(),
		aggregate.RecordedEvents(),
	)

	if err != nil {
		return xerrors.Errorf("customers[eventSourcedRepository].Persist: %w", err)
	}

	return nil
}
