package eventsourced

import (
	"database/sql"
	"go-iddd/service/customer/application"
	"go-iddd/service/lib/infrastructure/eventstore"
)

const streamPrefix = "customer"

type CustomersSessionStarter struct {
	eventStoreSessionStarter eventstore.StartsEventStoreSessions
}

func NewCustomersSessionStarter(
	eventStoreSessionStarter eventstore.StartsEventStoreSessions,
) *CustomersSessionStarter {

	return &CustomersSessionStarter{
		eventStoreSessionStarter: eventStoreSessionStarter,
	}
}

func (repo *CustomersSessionStarter) StartSession(tx *sql.Tx) application.Customers {
	return &Customers{
		eventStoreSession: repo.eventStoreSessionStarter.StartSession(tx),
	}
}
