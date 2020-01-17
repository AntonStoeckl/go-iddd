package service

import (
	"database/sql"
	"go-iddd/customer/application"
	"go-iddd/customer/infrastructure/eventsourced"
	"go-iddd/shared"
	"go-iddd/shared/infrastructure/eventstore"

	"github.com/cockroachdb/errors"
)

const (
	eventStoreTableName = "eventstore"
)

type DIContainer struct {
	postgresDBConn          *sql.DB
	unmarshalDomainEvent    shared.UnmarshalDomainEvent
	postgresEventStore      *eventstore.PostgresEventStore
	customersSessionStarter *eventsourced.CustomersSessionStarter
	customerCommandHandler  *application.CommandHandler
}

func NewDIContainer(
	postgresDBConn *sql.DB,
	unmarshalDomainEvent shared.UnmarshalDomainEvent,
) (*DIContainer, error) {

	if postgresDBConn == nil {
		return nil, errors.Mark(errors.New("newContainer: postgres DB connection must not be nil"), shared.ErrTechnical)
	}

	container := &DIContainer{
		postgresDBConn:       postgresDBConn,
		unmarshalDomainEvent: unmarshalDomainEvent,
	}

	return container, nil
}

func (container DIContainer) GetPostgresDBConn() *sql.DB {
	return container.postgresDBConn
}

func (container DIContainer) GetPostgresEventStore() *eventstore.PostgresEventStore {
	if container.postgresEventStore == nil {
		container.postgresEventStore = eventstore.NewPostgresEventStore(
			container.postgresDBConn,
			eventStoreTableName,
			container.unmarshalDomainEvent,
		)
	}

	return container.postgresEventStore
}

func (container DIContainer) GetCustomerRepository() *eventsourced.CustomersSessionStarter {
	if container.customersSessionStarter == nil {
		container.customersSessionStarter = eventsourced.NewCustomersSessionStarter(
			container.GetPostgresEventStore(),
		)
	}

	return container.customersSessionStarter
}

func (container DIContainer) GetCustomerCommandHandler() *application.CommandHandler {
	if container.customerCommandHandler == nil {
		container.customerCommandHandler = application.NewCommandHandler(
			container.GetCustomerRepository(),
			container.postgresDBConn,
		)
	}

	return container.customerCommandHandler
}
