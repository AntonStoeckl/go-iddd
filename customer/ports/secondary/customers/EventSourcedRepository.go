package customers

import (
	"database/sql"
	"go-iddd/customer/domain"
	"go-iddd/shared"
)

const streamPrefix = "customer"

type StartsEventStoreSessions interface {
	StartSession(tx *sql.Tx) shared.EventStore
}

type EventSourcedRepository struct {
	eventStoreSessionFactory StartsEventStoreSessions
	customerFactory          func(eventStream shared.DomainEvents) (domain.Customer, error)
	identityMap              *IdentityMap
}

func NewEventSourcedRepository(
	eventStoreSessionFactory StartsEventStoreSessions,
	customerFactory func(eventStream shared.DomainEvents) (domain.Customer, error),
	identityMap *IdentityMap,
) *EventSourcedRepository {

	return &EventSourcedRepository{
		eventStoreSessionFactory: eventStoreSessionFactory,
		customerFactory:          customerFactory,
		identityMap:              identityMap,
	}
}

func (repo *EventSourcedRepository) StartSession(tx *sql.Tx) *EventSourcedRepositorySession {
	return &EventSourcedRepositorySession{
		eventStoreSession: repo.eventStoreSessionFactory.StartSession(tx),
		customerFactory:   repo.customerFactory,
		identityMap:       repo.identityMap,
	}
}
