package customers

import (
	"database/sql"
	"go-iddd/customer/application"
	"go-iddd/customer/domain"
	"go-iddd/shared"
)

const streamPrefix = "customer"

type StartsEventStoreSessions interface {
	StartSession(tx *sql.Tx) shared.EventStore
}

type EventSourcedRepository struct {
	eventStoreSessionStarter StartsEventStoreSessions
	customerFactory          func(eventStream shared.DomainEvents) (*domain.Customer, error)
}

func NewEventSourcedRepository(
	eventStoreSessionStarter StartsEventStoreSessions,
	customerFactory func(eventStream shared.DomainEvents) (*domain.Customer, error),
) *EventSourcedRepository {

	return &EventSourcedRepository{
		eventStoreSessionStarter: eventStoreSessionStarter,
		customerFactory:          customerFactory,
	}
}

/***** Implement application.StartsCustomersSession *****/

func (repo *EventSourcedRepository) StartSession(tx *sql.Tx) application.Customers {
	return &EventSourcedRepositorySession{
		eventStoreSession: repo.eventStoreSessionStarter.StartSession(tx),
		customerFactory:   repo.customerFactory,
	}
}
