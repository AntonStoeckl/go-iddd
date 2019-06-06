package customers

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"sync"

	"golang.org/x/xerrors"
)

type customerFactory func(eventStream shared.DomainEvents) (domain.Customer, error)

type eventSourcedRepository struct {
	eventStore      shared.EventStore
	customerFactory customerFactory
	identityMap     map[string]domain.Customer
	identityMapMux  sync.Mutex
}

func NewEventSourcedRepository(
	eventStore shared.EventStore,
	customerFactory customerFactory,
) *eventSourcedRepository {

	return &eventSourcedRepository{
		eventStore:      eventStore,
		customerFactory: customerFactory,
		identityMap:     make(map[string]domain.Customer),
	}
}

/***** Implement domain.Customers *****/

func (repo *eventSourcedRepository) Register(customer domain.Customer) error {
	if _, found := repo.memorizedCustomerOf(customer.AggregateID().(*values.ID)); found {
		return xerrors.Errorf("customers[eventSourcedRepository].Register: already memorized in identityMap: %w", shared.ErrDuplicate)
	}

	if err := repo.eventStore.AppendToStream(customer.AggregateID(), customer.RecordedEvents()); err != nil {
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
		return memorizedCustomer, nil
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

/***** Methods for identityMap (local caching) *****/

func (repo *eventSourcedRepository) memorize(customer domain.Customer) {
	repo.identityMapMux.Lock()
	defer repo.identityMapMux.Unlock()

	repo.identityMap[customer.AggregateID().String()] = customer
}

func (repo *eventSourcedRepository) memorizedCustomerOf(id *values.ID) (domain.Customer, bool) {
	repo.identityMapMux.Lock()
	defer repo.identityMapMux.Unlock()

	customer, found := repo.identityMap[id.String()]

	if !found {
		return nil, false
	}

	return customer, true
}
