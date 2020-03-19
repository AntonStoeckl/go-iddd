package cmd

import (
	"database/sql"
	"go-iddd/service/customer/application"
	customercli "go-iddd/service/customer/infrastructure/primary/cli"
	customergrpc "go-iddd/service/customer/infrastructure/primary/grpc"
	"go-iddd/service/customer/infrastructure/secondary/eventstore"
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"go-iddd/service/lib/eventstore/postgres"

	"github.com/cockroachdb/errors"
)

const eventStoreTableName = "eventstore"

type DIContainer struct {
	postgresDBConn         *sql.DB
	unmarshalCustomerEvent es.UnmarshalDomainEvent
	customerEventStore     *eventstore.CustomerEventStore
	customerCommandHandler *application.CustomerCommandHandler
	customerQueryHandler   *application.CustomerQueryHandler
	customerServer         customergrpc.CustomerServer
	customerApp            *customercli.CustomerApp
}

func NewDIContainer(
	postgresDBConn *sql.DB,
	unmarshalDomainEventForWriteModel es.UnmarshalDomainEvent,
) (*DIContainer, error) {

	if postgresDBConn == nil {
		return nil, errors.Mark(errors.New("newDIContainer: postgres DB connection must not be nil"), lib.ErrTechnical)
	}

	container := &DIContainer{
		postgresDBConn:         postgresDBConn,
		unmarshalCustomerEvent: unmarshalDomainEventForWriteModel,
	}

	container.init()

	return container, nil
}

func (container DIContainer) init() {
	container.GetCustomerEventStore()
	container.GetCustomerCommandHandler()
	container.GetCustomerQueryHandler()
	container.GetCustomerServer()
	container.GetCustomerApp()
}

func (container DIContainer) GetPostgresDBConn() *sql.DB {
	return container.postgresDBConn
}

func (container DIContainer) GetCustomerEventStore() *eventstore.CustomerEventStore {
	if container.customerEventStore == nil {
		container.customerEventStore = eventstore.NewCustomerEventStore(
			postgres.NewEventStore(
				container.postgresDBConn,
				eventStoreTableName,
				container.unmarshalCustomerEvent,
			),
		)
	}

	return container.customerEventStore
}

func (container DIContainer) GetCustomerCommandHandler() *application.CustomerCommandHandler {
	if container.customerCommandHandler == nil {
		container.customerCommandHandler = application.NewCustomerCommandHandler(
			container.GetCustomerEventStore(),
		)
	}

	return container.customerCommandHandler
}

func (container DIContainer) GetCustomerQueryHandler() *application.CustomerQueryHandler {
	if container.customerQueryHandler == nil {
		container.customerQueryHandler = application.NewCustomerQueryHandler(
			container.GetCustomerEventStore(),
		)
	}

	return container.customerQueryHandler
}

func (container DIContainer) GetCustomerServer() customergrpc.CustomerServer {
	if container.customerServer == nil {
		container.customerServer = customergrpc.NewCustomerServer(
			container.GetCustomerCommandHandler().RegisterCustomer,
			container.GetCustomerCommandHandler().ConfirmCustomerEmailAddress,
			container.GetCustomerCommandHandler().ChangeCustomerEmailAddress,
		)
	}

	return container.customerServer
}

func (container DIContainer) GetCustomerApp() *customercli.CustomerApp {
	if container.customerApp == nil {
		container.customerApp = customercli.NewCustomerApp(
			container.GetCustomerCommandHandler().RegisterCustomer,
			container.GetCustomerCommandHandler().ConfirmCustomerEmailAddress,
			container.GetCustomerCommandHandler().ChangeCustomerEmailAddress,
		)
	}

	return container.customerApp
}
