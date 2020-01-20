package cmd

import (
	"database/sql"
	"go-iddd/service/customer/application"
	"go-iddd/service/customer/infrastructure/eventsourced"
	"go-iddd/service/lib"
	"go-iddd/service/lib/infrastructure/eventstore"

	"github.com/cockroachdb/errors"
)

const (
	eventStoreTableName = "eventstore"
)

type DIContainer struct {
	postgresDBConn          *sql.DB
	unmarshalDomainEvent    lib.UnmarshalDomainEvent
	postgresEventStore      *eventstore.PostgresEventStore
	customersSessionStarter *eventsourced.CustomersSessionStarter
	customerCommandHandler  *application.CommandHandler
}

func NewDIContainer(
	postgresDBConn *sql.DB,
	unmarshalDomainEvent lib.UnmarshalDomainEvent,
) (*DIContainer, error) {

	if postgresDBConn == nil {
		return nil, errors.Mark(errors.New("newContainer: postgres DB connection must not be nil"), lib.ErrTechnical)
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
