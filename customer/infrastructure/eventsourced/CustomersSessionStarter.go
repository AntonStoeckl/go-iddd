package eventsourced

import (
	"database/sql"
	"go-iddd/customer/application"
	"go-iddd/customer/domain"
	"go-iddd/shared"
	"go-iddd/shared/infrastructure/eventstore"
)

const streamPrefix = "customer"

type CustomersSessionStarter struct {
	eventStoreSessionStarter eventstore.StartsEventStoreSessions
	customerFactory          func(eventStream shared.DomainEvents) (*domain.Customer, error)
}

func NewCustomersSessionStarter(
	eventStoreSessionStarter eventstore.StartsEventStoreSessions,
	customerFactory func(eventStream shared.DomainEvents) (*domain.Customer, error),
) *CustomersSessionStarter {

	return &CustomersSessionStarter{
		eventStoreSessionStarter: eventStoreSessionStarter,
		customerFactory:          customerFactory,
	}
}

func (repo *CustomersSessionStarter) StartSession(tx *sql.Tx) application.Customers {
	return &Customers{
		eventStoreSession: repo.eventStoreSessionStarter.StartSession(tx),
		customerFactory:   repo.customerFactory,
	}
}
