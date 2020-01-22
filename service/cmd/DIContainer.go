package cmd

import (
	"database/sql"
	"go-iddd/service/customer/application"
	customercli "go-iddd/service/customer/infrastructure/primary/cli"
	customergrpc "go-iddd/service/customer/infrastructure/primary/grpc"
	"go-iddd/service/customer/infrastructure/secondary/eventsourced"
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
	customerServer          customergrpc.CustomerServer
	customerApp             *customercli.CustomerApp
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

	container.init()

	return container, nil
}

func (container DIContainer) init() {
	container.GetPostgresEventStore()
	container.GetCustomerRepository()
	container.GetCustomerCommandHandler()
	container.GetCustomerServer()
	container.GetCustomerApp()
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
