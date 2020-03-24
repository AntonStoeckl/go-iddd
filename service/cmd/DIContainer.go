package cmd

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/service/customer/application/command"
	"github.com/AntonStoeckl/go-iddd/service/customer/application/query"
	customercli "github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/primary/cli"
	customergrpc "github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/primary/grpc"
	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/secondary/eventstore"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/AntonStoeckl/go-iddd/service/lib/eventstore/postgres"
	"github.com/cockroachdb/errors"
)

const eventStoreTableName = "eventstore"

type DIContainer struct {
	postgresDBConn         *sql.DB
	unmarshalCustomerEvent es.UnmarshalDomainEvent
	customerEventStore     *eventstore.CustomerEventStore
	customerCommandHandler *command.CustomerCommandHandler
	customerQueryHandler   *query.CustomerQueryHandler
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

func (container DIContainer) GetCustomerCommandHandler() *command.CustomerCommandHandler {
	if container.customerCommandHandler == nil {
		container.customerCommandHandler = command.NewCustomerCommandHandler(
			container.GetCustomerEventStore(),
		)
	}

	return container.customerCommandHandler
}

func (container DIContainer) GetCustomerQueryHandler() *query.CustomerQueryHandler {
	if container.customerQueryHandler == nil {
		container.customerQueryHandler = query.NewCustomerQueryHandler(
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
			container.GetCustomerCommandHandler().ChangeCustomerName,
			container.GetCustomerCommandHandler().DeleteCustomer,
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
