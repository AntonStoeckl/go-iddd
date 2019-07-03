package customers

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"math"
	"sync"

	"golang.org/x/xerrors"
)

type customerFactory func(eventStream shared.DomainEvents) (domain.Customer, error)

type eventSourcedRepository struct {
	eventStore      shared.EventStore
	customerFactory customerFactory
	cache           map[string]shared.EventsourcedAggregate
	cacheMux        sync.Mutex
}

func NewEventSourcedRepository(
	eventStore shared.EventStore,
	customerFactory customerFactory,
) *eventSourcedRepository {

	return &eventSourcedRepository{
		eventStore:      eventStore,
		customerFactory: customerFactory,
		cache:           make(map[string]shared.EventsourcedAggregate),
	}
}

/***** Implement domain.Customers *****/

func (repo *eventSourcedRepository) Register(customer domain.Customer) error {
	if _, found := repo.memorizedCustomerOf(customer.AggregateID().(*values.ID)); found {
		return xerrors.Errorf("customers[eventSourcedRepository].Register: already memorized in cache: %w", shared.ErrDuplicate)
	}

	if err := repo.eventStore.AppendToStream(customer.RecordedEvents(true)); err != nil {
		if xerrors.Is(err, shared.ErrConcurrencyConflict) {
			return xerrors.Errorf("customers[eventSourcedRepository].Register: %s: %w", err, shared.ErrDuplicate)
		}

		return xerrors.Errorf("customers[eventSourcedRepository].Register: %w", err)
	}

	repo.memorize(customer)

	return nil
}

func (repo *eventSourcedRepository) Of(id *values.ID) (domain.Customer, error) {
	if memorizedCustomer, found := repo.memorizedCustomerOf(id); found {
		latestEvents, err := repo.eventStore.LoadPartialEventStream(id, memorizedCustomer.StreamVersion(), uint(math.MaxUint32))
		if err != nil {
			return nil, xerrors.Errorf("customers[eventSourcedRepository].Of: %w", err)
		}

		memorizedCustomer.Apply(latestEvents)

		return memorizedCustomer.(domain.Customer).Clone(), nil
	}

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

	repo.memorize(customer)

	return customer.Clone(), nil
}

/***** Implement shared.PersistsEventsourcedAggregates *****/

func (repo *eventSourcedRepository) Persist(aggregate shared.EventsourcedAggregate) error {
	err := repo.eventStore.AppendToStream(aggregate.RecordedEvents(true))

	if err != nil {
		return xerrors.Errorf("customers[eventSourcedRepository].Persist: %w", err)
	}

	repo.memorize(aggregate)

	return nil
}

/***** Methods for local caching *****/

func (repo *eventSourcedRepository) memorize(aggregate shared.EventsourcedAggregate) {
	repo.cacheMux.Lock()
	defer repo.cacheMux.Unlock()

	repo.cache[aggregate.AggregateID().String()] = aggregate
}

func (repo *eventSourcedRepository) memorizedCustomerOf(id *values.ID) (shared.EventsourcedAggregate, bool) {
	repo.cacheMux.Lock()
	defer repo.cacheMux.Unlock()

	aggregate, found := repo.cache[id.String()]

	if !found {
		return nil, false
	}

	return aggregate, true
}
