package customers

import (
	"go-iddd/customer/domain"
	"go-iddd/shared"

	"golang.org/x/xerrors"
)

const streamPrefix = "customer"

type EventSourcedRepository struct {
	eventStore      shared.StartsEventStoreSessions
	customerFactory func(eventStream shared.DomainEvents) (domain.Customer, error)
	identityMap     *IdentityMap
}

func NewEventSourcedRepository(
	eventStore shared.StartsEventStoreSessions,
	customerFactory func(eventStream shared.DomainEvents) (domain.Customer, error),
	identityMap *IdentityMap,
) *EventSourcedRepository {

	return &EventSourcedRepository{
		eventStore:      eventStore,
		customerFactory: customerFactory,
		identityMap:     identityMap,
	}
}

func (repo *EventSourcedRepository) StartSession() (*EventSourcedRepositorySession, error) {
	eventStoreSession, err := repo.eventStore.StartSession()
	if err != nil {
		return nil, xerrors.Errorf("eventSourcedRepository.StartSession: %w", err)
	}

	repoSession := &EventSourcedRepositorySession{
		eventStoreSession: eventStoreSession,
		customerFactory:   repo.customerFactory,
		identityMap:       repo.identityMap,
	}

	return repoSession, nil
}
