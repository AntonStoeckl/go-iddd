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

func UsePostgresDBConn(postgresDBConn *sql.DB) DIOption {
	return func(container *DIContainer) error {
		if postgresDBConn == nil {
			return errors.New("pgDBConn must not be nil")
		}

		container.infra.pgDBConn = postgresDBConn

		return nil
	}
}

func ReplaceGRPCCustomerServer(customerGRPCServer customergrpc.CustomerServer) DIOption {
	return func(container *DIContainer) error {
		if customerGRPCServer == nil {
			return errors.New("grpcCustomerServer must not be nil")
		}

		container.service.grpcCustomerServer = customerGRPCServer

		return nil
	}
}

type DIContainer struct {
	config *Config

	infra struct {
		pgDBConn *sql.DB
	}

	dependency struct {
		marshalCustomerEvent              es.MarshalDomainEvent
		unmarshalCustomerEvent            es.UnmarshalDomainEvent
		buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions
	}

	service struct {
		customerEventStore     *postgres.CustomerEventStore
		customerCommandHandler *application.CustomerCommandHandler
		customerQueryHandler   *application.CustomerQueryHandler
		grpcCustomerServer     customergrpc.CustomerServer
		grpcServer             *grpc.Server
	}
}

func MustBuildDIContainer(
	config *Config,
	logger *shared.Logger,
	marshalCustomerEvent es.MarshalDomainEvent,
	unmarshalCustomerEvent es.UnmarshalDomainEvent,
	buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions,
	opts ...DIOption,
) *DIContainer {

	container := &DIContainer{}
	container.config = config
	container.dependency.marshalCustomerEvent = marshalCustomerEvent
	container.dependency.unmarshalCustomerEvent = unmarshalCustomerEvent
	container.dependency.buildUniqueEmailAddressAssertions = buildUniqueEmailAddressAssertions

	for _, opt := range opts {
		if err := opt(container); err != nil {
			logger.Panicf("mustBuildDIContainer: %s", err)
		}
	}

	return container
}

func (container DIContainer) GetCustomerEventStore() *postgres.CustomerEventStore {
	if container.service.customerEventStore == nil {
		container.service.customerEventStore = postgres.NewCustomerEventStore(
			container.infra.pgDBConn,
			eventStoreTableName,
			container.dependency.marshalCustomerEvent,
			container.dependency.unmarshalCustomerEvent,
			uniqueEmailAddressesTableName,
			container.dependency.buildUniqueEmailAddressAssertions,
		)
	}

	return container.service.customerEventStore
}

func (container DIContainer) GetCustomerCommandHandler() *application.CustomerCommandHandler {
	if container.service.customerCommandHandler == nil {
		container.service.customerCommandHandler = application.NewCustomerCommandHandler(
			container.GetCustomerEventStore().RetrieveEventStream,
			container.GetCustomerEventStore().StartEventStream,
			container.GetCustomerEventStore().AppendToEventStream,
		)
	}

	return container.service.customerCommandHandler
}

func (container DIContainer) GetCustomerQueryHandler() *application.CustomerQueryHandler {
	if container.service.customerQueryHandler == nil {
		container.service.customerQueryHandler = application.NewCustomerQueryHandler(
			container.GetCustomerEventStore().RetrieveEventStream,
		)
	}

	return container.service.customerQueryHandler
}

func (container DIContainer) GetGRPCCustomerServer() customergrpc.CustomerServer {
	if container.service.grpcCustomerServer == nil {
		container.service.grpcCustomerServer = customergrpc.NewCustomerServer(
			container.GetCustomerCommandHandler().RegisterCustomer,
			container.GetCustomerCommandHandler().ConfirmCustomerEmailAddress,
			container.GetCustomerCommandHandler().ChangeCustomerEmailAddress,
			container.GetCustomerCommandHandler().ChangeCustomerName,
			container.GetCustomerCommandHandler().DeleteCustomer,
			container.GetCustomerQueryHandler().CustomerViewByID,
		)
	}

	return container.service.grpcCustomerServer
}

func (container DIContainer) GetGRPCServer() *grpc.Server {
	if container.service.grpcServer == nil {
		container.service.grpcServer = grpc.NewServer()
		customergrpc.RegisterCustomerServer(container.service.grpcServer, container.GetGRPCCustomerServer())
		reflection.Register(container.service.grpcServer)
	}

	return container.service.grpcServer
}
