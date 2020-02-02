package cmd

import (
	"database/sql"
	"go-iddd/service/customer/application"
	customercli "go-iddd/service/customer/infrastructure/primary/cli"
	customergrpc "go-iddd/service/customer/infrastructure/primary/grpc"
	"go-iddd/service/customer/infrastructure/secondary/forstoringcustomerevents/eventstore"
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"go-iddd/service/lib/eventstore/postgres"

	"github.com/cockroachdb/errors"
)

const eventStoreTableName = "eventstore"

type DIContainer struct {
	postgresDBConn         *sql.DB
	unmarshalDomainEvent   es.UnmarshalDomainEvent
	eventStore             *postgres.EventStore
	customerEventStore     *eventstore.CustomerEventStore
	customerCommandHandler *application.CommandHandler
	customerServer         customergrpc.CustomerServer
	customerApp            *customercli.CustomerApp
}

func NewDIContainer(
	postgresDBConn *sql.DB,
	unmarshalDomainEvent es.UnmarshalDomainEvent,
) (*DIContainer, error) {

	if postgresDBConn == nil {
		return nil, errors.Mark(errors.New("newContainer: postgres DB connection must not be nil"), lib.ErrTechnical)
	}

	container := &DIContainer{
		postgresDBConn:       postgresDBConn,
		unmarshalDomainEvent: unmarshalDomainEvent,
	}

	container.init()

	return container, nil
}

func (container DIContainer) init() {
	container.getEventStore()
	container.GetCustomerEventStore()
	container.GetCustomerCommandHandler()
	container.GetCustomerServer()
	container.GetCustomerApp()
}

func (container DIContainer) GetPostgresDBConn() *sql.DB {
	return container.postgresDBConn
}

func (container DIContainer) getEventStore() *postgres.EventStore {
	if container.eventStore == nil {
		container.eventStore = postgres.NewEventStore(
			container.postgresDBConn,
			eventStoreTableName,
			container.unmarshalDomainEvent,
		)
	}

	return container.eventStore
}

func (container DIContainer) GetCustomerEventStore() *eventstore.CustomerEventStore {
	if container.customerEventStore == nil {
		container.customerEventStore = eventstore.NewCustomerEventStore(
			container.getEventStore(),
		)
	}

	return container.customerEventStore
}

func (container DIContainer) GetCustomerCommandHandler() *application.CommandHandler {
	if container.customerCommandHandler == nil {
		container.customerCommandHandler = application.NewCommandHandler(
			container.GetCustomerEventStore(),
		)
	}

	return container.customerCommandHandler
}

func (container DIContainer) GetCustomerServer() customergrpc.CustomerServer {
	if container.customerServer == nil {
		container.customerServer = customergrpc.NewCustomerServer(
			container.GetCustomerCommandHandler(),
			container.GetCustomerCommandHandler(),
			container.GetCustomerCommandHandler(),
		)
	}

	return container.customerServer
}

func (container DIContainer) GetCustomerApp() *customercli.CustomerApp {
	if container.customerApp == nil {
		container.customerApp = customercli.NewCustomerApp(
			container.GetCustomerCommandHandler(),
			container.GetCustomerCommandHandler(),
			container.GetCustomerCommandHandler(),
		)
	}

	return container.customerApp
}
