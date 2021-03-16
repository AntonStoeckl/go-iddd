package grpc

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer"
	customergrpc "github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/grpc"
	customergrpcproto "github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/grpc/proto"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/postgres"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/serialization"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	eventStoreTableName           = "eventstore"
	uniqueEmailAddressesTableName = "unique_email_addresses"
	uniqueIdentitiesTableName     = "unique_identities"
	maxConcurrencyConflictRetries = 10
)

type DIContainer struct {
	config *Config

	infra struct {
		pgDBConn *sql.DB
	}

	dependency struct {
		// Customer
		marshalCustomerEvent              es.MarshalDomainEvent
		unmarshalCustomerEvent            es.UnmarshalDomainEvent
		buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions

		// Identity
		marshalIdentityEvent   es.MarshalDomainEvent
		unmarshalIdentityEvent es.UnmarshalDomainEvent
	}

	service struct {
		// Customer
		customerEventStore     *postgres.CustomerEventStore
		customerCommandHandler *application.CustomerCommandHandler
		customerQueryHandler   *application.CustomerQueryHandler
		grpcCustomerServer     customergrpcproto.CustomerServer

		// Identity
		uniqueIdentities       *postgres.UniqueIdentities
		identityEventStore     *postgres.IdentityEventStore
		identityCommandHandler *application.IdentityCommandHandler
		loginHandler           *application.LoginHandler

		// Generic
		eventStore *es.EventStore
		grpcServer *grpc.Server
	}
}

// MustBuildDIContainer - factory method to build a DIContainer with DIOptions
//   Panics if it fails to apply an Option
func MustBuildDIContainer(config *Config, logger *shared.Logger, opts ...DIOption) *DIContainer {
	container := &DIContainer{}
	container.config = config

	/*** Define default dependencies ***/
	container.dependency.marshalCustomerEvent = serialization.MarshalCustomerEvent
	container.dependency.unmarshalCustomerEvent = serialization.UnmarshalCustomerEvent
	container.dependency.buildUniqueEmailAddressAssertions = customer.BuildUniqueEmailAddressAssertions
	container.dependency.marshalIdentityEvent = serialization.MarshalIdentityEvent
	container.dependency.unmarshalIdentityEvent = serialization.UnmarshalIdentityEvent

	/*** Apply options for infra, dependencies, services ***/
	for _, opt := range opts {
		if err := opt(container); err != nil {
			logger.Panic().Msgf("mustBuildDIContainer: %s", err)
		}
	}

	container.init()

	return container
}

// init - initializes all dependencies in advance so we have no lazy initialization
func (container *DIContainer) init() {
	// Customer
	_ = container.GetCustomerEventStore()
	_ = container.GetCustomerCommandHandler()
	_ = container.GetCustomerQueryHandler()
	_ = container.getGRPCCustomerServer()

	// Identity
	_ = container.GetUniqueIdentities()
	_ = container.GetIdentityEventStore()
	_ = container.GetIdentityCommandHandler()
	_ = container.GetLoginHandler()

	// Generic
	_ = container.getEventStore()
	_ = container.GetGRPCServer()
}

/*
 * Customer
 */

func (container *DIContainer) GetCustomerEventStore() *postgres.CustomerEventStore {
	if container.service.customerEventStore == nil {
		uniqueCustomerEmailAddresses := postgres.NewUniqueCustomerEmailAddresses(
			uniqueEmailAddressesTableName,
			container.dependency.buildUniqueEmailAddressAssertions,
		)

		container.service.customerEventStore = postgres.NewCustomerEventStore(
			container.infra.pgDBConn,
			container.getEventStore().RetrieveEventStream,
			container.getEventStore().AppendEventsToStream,
			container.getEventStore().PurgeEventStream,
			uniqueCustomerEmailAddresses.AssertUniqueEmailAddress,
			uniqueCustomerEmailAddresses.PurgeUniqueEmailAddress,
			container.dependency.marshalCustomerEvent,
			container.dependency.unmarshalCustomerEvent,
		)
	}

	return container.service.customerEventStore
}

func (container *DIContainer) GetCustomerCommandHandler() *application.CustomerCommandHandler {
	if container.service.customerCommandHandler == nil {
		container.service.customerCommandHandler = application.NewCustomerCommandHandler(
			container.GetCustomerEventStore().RetrieveEventStream,
			container.GetCustomerEventStore().StartEventStream,
			container.GetCustomerEventStore().AppendToEventStream,
		)
	}

	return container.service.customerCommandHandler
}

func (container *DIContainer) GetCustomerQueryHandler() *application.CustomerQueryHandler {
	if container.service.customerQueryHandler == nil {
		container.service.customerQueryHandler = application.NewCustomerQueryHandler(
			container.GetCustomerEventStore().RetrieveEventStream,
		)
	}

	return container.service.customerQueryHandler
}

func (container *DIContainer) getGRPCCustomerServer() customergrpcproto.CustomerServer {
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

/*
 * Identity
 */

func (container *DIContainer) GetUniqueIdentities() *postgres.UniqueIdentities {
	if container.service.uniqueIdentities == nil {
		container.service.uniqueIdentities = postgres.NewUniqueIdentities(
			uniqueIdentitiesTableName,
			container.infra.pgDBConn,
		)
	}

	return container.service.uniqueIdentities
}

func (container *DIContainer) GetIdentityEventStore() *postgres.IdentityEventStore {
	if container.service.identityEventStore == nil {
		container.service.identityEventStore = postgres.NewIdentityEventStore(
			container.infra.pgDBConn,
			container.getEventStore().RetrieveEventStream,
			container.getEventStore().AppendEventsToStream,
			container.getEventStore().PurgeEventStream,
			container.dependency.marshalIdentityEvent,
			container.dependency.unmarshalIdentityEvent,
		)
	}

	return container.service.identityEventStore
}

func (container *DIContainer) GetIdentityCommandHandler() *application.IdentityCommandHandler {
	if container.service.identityCommandHandler == nil {
		container.service.identityCommandHandler = application.NewIdentityCommandHandler(
			container.infra.pgDBConn,
			container.GetUniqueIdentities(),
			container.GetIdentityEventStore(),
			maxConcurrencyConflictRetries,
		)
	}

	return container.service.identityCommandHandler
}

func (container *DIContainer) GetLoginHandler() *application.LoginHandler {
	if container.service.loginHandler == nil {
		container.service.loginHandler = application.NewLoginHandler(
			container.GetUniqueIdentities(),
			container.GetIdentityEventStore(),
		)
	}

	return container.service.loginHandler
}

/*
 * Generic
 */

func (container *DIContainer) GetPostgresDBConn() *sql.DB {
	return container.infra.pgDBConn
}

func (container *DIContainer) getEventStore() *es.EventStore {
	if container.service.eventStore == nil {
		container.service.eventStore = es.NewEventStore(eventStoreTableName)
	}

	return container.service.eventStore
}

func (container *DIContainer) GetGRPCServer() *grpc.Server {
	if container.service.grpcServer == nil {
		container.service.grpcServer = grpc.NewServer()
		customergrpcproto.RegisterCustomerServer(container.service.grpcServer, container.getGRPCCustomerServer())
		reflection.Register(container.service.grpcServer)
	}

	return container.service.grpcServer
}

/*
 * Options
 */

type DIOption func(container *DIContainer) error

func UsePostgresDBConn(dbConn *sql.DB) DIOption {
	return func(container *DIContainer) error {
		if dbConn == nil {
			return errors.New("pgDBConn must not be nil")
		}

		container.infra.pgDBConn = dbConn

		return nil
	}
}

func WithMarshalCustomerEvents(fn es.MarshalDomainEvent) DIOption {
	return func(container *DIContainer) error {
		container.dependency.marshalCustomerEvent = fn
		return nil
	}
}

func WithUnmarshalCustomerEvents(fn es.UnmarshalDomainEvent) DIOption {
	return func(container *DIContainer) error {
		container.dependency.unmarshalCustomerEvent = fn
		return nil
	}
}

func WithMarshalIdentityEvents(fn es.MarshalDomainEvent) DIOption {
	return func(container *DIContainer) error {
		container.dependency.marshalIdentityEvent = fn
		return nil
	}
}

func WithUnmarshalIdentityEvents(fn es.UnmarshalDomainEvent) DIOption {
	return func(container *DIContainer) error {
		container.dependency.unmarshalIdentityEvent = fn
		return nil
	}
}

func WithBuildUniqueEmailAddressAssertions(fn customer.ForBuildingUniqueEmailAddressAssertions) DIOption {
	return func(container *DIContainer) error {
		container.dependency.buildUniqueEmailAddressAssertions = fn
		return nil
	}
}

func ReplaceGRPCCustomerServer(server customergrpcproto.CustomerServer) DIOption {
	return func(container *DIContainer) error {
		if server == nil {
			return errors.New("grpcCustomerServer must not be nil")
		}

		container.service.grpcCustomerServer = server

		return nil
	}
}
