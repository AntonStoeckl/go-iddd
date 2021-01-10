package cmd

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer"
	customergrpc "github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/adapter/grpc"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/adapter/postgres"
	"github.com/AntonStoeckl/go-iddd/service/shared"
	"github.com/AntonStoeckl/go-iddd/service/shared/es"
	"github.com/cockroachdb/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	eventStoreTableName           = "eventstore"
	uniqueEmailAddressesTableName = "unique_email_addresses"
)

type DIOption func(container *DIContainer) error

func WithPostgresDBConn(postgresDBConn *sql.DB) DIOption {
	return func(container *DIContainer) error {
		if postgresDBConn == nil {
			return errors.New("postgresDBConn must not be nil")
		}

		container.postgresDBConn = postgresDBConn

		return nil
	}
}

func WithGRPCCustomerServer(customerGRPCServer customergrpc.CustomerServer) DIOption {
	return func(container *DIContainer) error {
		if customerGRPCServer == nil {
			return errors.New("grpcCustomerServer must not be nil")
		}

		container.grpcCustomerServer = customerGRPCServer

		return nil
	}
}

type DIContainer struct {
	postgresDBConn                    *sql.DB
	customerEventStore                *postgres.CustomerEventStore
	marshalCustomerEvent              es.MarshalDomainEvent
	unmarshalCustomerEvent            es.UnmarshalDomainEvent
	buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions
	customerCommandHandler            *application.CustomerCommandHandler
	customerQueryHandler              *application.CustomerQueryHandler
	grpcCustomerServer                customergrpc.CustomerServer
	grpcServer                        *grpc.Server
}

func NewDIContainer(
	marshalCustomerEvent es.MarshalDomainEvent,
	unmarshalCustomerEvent es.UnmarshalDomainEvent,
	buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions,
	opts ...DIOption,
) (*DIContainer, error) {

	container := &DIContainer{
		marshalCustomerEvent:              marshalCustomerEvent,
		unmarshalCustomerEvent:            unmarshalCustomerEvent,
		buildUniqueEmailAddressAssertions: buildUniqueEmailAddressAssertions,
	}

	for _, opt := range opts {
		if err := opt(container); err != nil {
			return nil, shared.MarkAndWrapError(err, shared.ErrTechnical, "newDIContainer")
		}
	}

	container.init()

	return container, nil
}

func (container DIContainer) init() {
	_ = container.GetCustomerEventStore()
	_ = container.GetCustomerCommandHandler()
	_ = container.GetCustomerQueryHandler()
	_ = container.GetGRPCCustomerServer()
	_ = container.GetGRPCServer()
}

func (container DIContainer) GetPostgresDBConn() *sql.DB {
	return container.postgresDBConn
}

func (container DIContainer) GetCustomerEventStore() *postgres.CustomerEventStore {
	if container.customerEventStore == nil {
		container.customerEventStore = postgres.NewCustomerEventStore(
			container.postgresDBConn,
			eventStoreTableName,
			container.marshalCustomerEvent,
			container.unmarshalCustomerEvent,
			uniqueEmailAddressesTableName,
			container.buildUniqueEmailAddressAssertions,
		)
	}

	return container.customerEventStore
}

func (container DIContainer) GetCustomerCommandHandler() *application.CustomerCommandHandler {
	if container.customerCommandHandler == nil {
		container.customerCommandHandler = application.NewCustomerCommandHandler(
			container.GetCustomerEventStore().RetrieveEventStream,
			container.GetCustomerEventStore().StartEventStream,
			container.GetCustomerEventStore().AppendToEventStream,
		)
	}

	return container.customerCommandHandler
}

func (container DIContainer) GetCustomerQueryHandler() *application.CustomerQueryHandler {
	if container.customerQueryHandler == nil {
		container.customerQueryHandler = application.NewCustomerQueryHandler(
			container.GetCustomerEventStore().RetrieveEventStream,
		)
	}

	return container.customerQueryHandler
}

func (container DIContainer) GetGRPCCustomerServer() customergrpc.CustomerServer {
	if container.grpcCustomerServer == nil {
		container.grpcCustomerServer = customergrpc.NewCustomerServer(
			container.GetCustomerCommandHandler().RegisterCustomer,
			container.GetCustomerCommandHandler().ConfirmCustomerEmailAddress,
			container.GetCustomerCommandHandler().ChangeCustomerEmailAddress,
			container.GetCustomerCommandHandler().ChangeCustomerName,
			container.GetCustomerCommandHandler().DeleteCustomer,
			container.GetCustomerQueryHandler().CustomerViewByID,
		)
	}

	return container.grpcCustomerServer
}

func (container DIContainer) GetGRPCServer() *grpc.Server {
	if container.grpcServer == nil {
		container.grpcServer = grpc.NewServer()
		customergrpc.RegisterCustomerServer(container.grpcServer, container.GetGRPCCustomerServer())
		reflection.Register(container.grpcServer)
	}

	return container.grpcServer
}
