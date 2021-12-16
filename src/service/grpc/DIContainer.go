package grpc

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer"
	customergrpc "github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/grpc"
	customergrpcproto "github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/grpc/proto"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/mongodb"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/postgres"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/serialization"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	eventStoreTableName           = "eventstore"
	uniqueEmailAddressesTableName = "unique_email_addresses"
)

type DIOption func(container *DIContainer) error

func UsePostgresDBConn(dbConn *sql.DB) DIOption {
	return func(container *DIContainer) error {
		if (dbConn == nil) && (container.config.EventStoreDB == "postgres") {
			return errors.New("pgDBConn must not be nil")
		}

		container.infra.pgDBConn = dbConn

		return nil
	}
}
func UseMongoDBConn(dbConn *mongo.Client) DIOption {
	return func(container *DIContainer) error {
		if (dbConn == nil) && (container.config.EventStoreDB == "mongodb") {
			return errors.New("mongodbConn must not be nil")
		}

		container.infra.mongodbConn = dbConn

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

type DIContainer struct {
	config *Config

	infra struct {
		pgDBConn    *sql.DB
		mongodbConn *mongo.Client
	}

	dependency struct {
		marshalCustomerEvent              es.MarshalDomainEvent
		unmarshalCustomerEvent            es.UnmarshalDomainEvent
		buildUniqueEmailAddressAssertions customer.ForBuildingUniqueEmailAddressAssertions
	}

	service struct {
		customerEventStore     application.EventStoreInterface
		customerCommandHandler *application.CustomerCommandHandler
		customerQueryHandler   *application.CustomerQueryHandler
		grpcCustomerServer     customergrpcproto.CustomerServer
		grpcServer             *grpc.Server
	}
}

func MustBuildDIContainer(config *Config, logger *shared.Logger, opts ...DIOption) *DIContainer {
	container := &DIContainer{}
	container.config = config

	/*** Define default dependencies ***/
	container.dependency.marshalCustomerEvent = serialization.MarshalCustomerEvent
	container.dependency.unmarshalCustomerEvent = serialization.UnmarshalCustomerEvent
	container.dependency.buildUniqueEmailAddressAssertions = customer.BuildUniqueEmailAddressAssertions

	/*** Apply options for infra, dependencies, services ***/
	for _, opt := range opts {
		if err := opt(container); err != nil {
			logger.Panic().Msgf("mustBuildDIContainer: %s", err)
		}
	}

	container.init()

	return container
}

func (container *DIContainer) init() {
	_ = container.GetCustomerEventStore()
	_ = container.GetCustomerCommandHandler()
	_ = container.GetCustomerQueryHandler()
	_ = container.getGRPCCustomerServer()
	_ = container.GetGRPCServer()
}

func (container *DIContainer) GetPostgresDBConn() *sql.DB {
	return container.infra.pgDBConn
}
func (container *DIContainer) GetMongoDBConn() *mongo.Client {
	return container.infra.mongodbConn
}

func (container *DIContainer) getPostgresEventStore() *es.PostgresEventStore {

	return es.NewPostgresEventStore(
		eventStoreTableName,
		container.dependency.marshalCustomerEvent,
		container.dependency.unmarshalCustomerEvent,
	)
}
func (container *DIContainer) getMongoEventStore() *es.MongodbEventStore {
	collection := container.infra.mongodbConn.Database(container.config.Mongodb.MongoInitdbDatabase).Collection(eventStoreTableName)
	return es.NewMongodbEventStore(
		collection,
		container.dependency.marshalCustomerEvent,
		container.dependency.unmarshalCustomerEvent,
	)
}

func (container *DIContainer) GetCustomerEventStore() application.EventStoreInterface {

	switch container.config.EventStoreDB {
	case "postgres":
		if container.service.customerEventStore == nil {
			uniqueCustomerEmailAddresses := postgres.NewUniqueCustomerEmailAddresses(
				uniqueEmailAddressesTableName,
				container.dependency.buildUniqueEmailAddressAssertions,
			)

			container.service.customerEventStore = postgres.NewCustomerPostgresEventStore(
				container.infra.pgDBConn,
				container.getPostgresEventStore().RetrieveEventStream,
				container.getPostgresEventStore().AppendEventsToStream,
				container.getPostgresEventStore().PurgeEventStream,
				uniqueCustomerEmailAddresses.AssertUniqueEmailAddress,
				uniqueCustomerEmailAddresses.PurgeUniqueEmailAddress,
			)
		}
	case "mongodb":
		if container.service.customerEventStore == nil {
			emailCollection := container.infra.mongodbConn.Database(container.config.Mongodb.MongoInitdbDatabase).Collection(uniqueEmailAddressesTableName)
			uniqueCustomerEmailAddresses := mongodb.NewUniqueCustomerEmailAddresses(
				emailCollection,
				container.dependency.buildUniqueEmailAddressAssertions,
			)
			collection := container.infra.mongodbConn.Database(container.config.Mongodb.MongoInitdbDatabase).Collection(eventStoreTableName)
			container.service.customerEventStore = mongodb.NewCustomerMongodbEventStore(
				collection,
				container.getMongoEventStore().RetrieveEventStream,
				container.getMongoEventStore().AppendEventsToStream,
				container.getMongoEventStore().PurgeEventStream,
				uniqueCustomerEmailAddresses.AssertUniqueEmailAddress,
				uniqueCustomerEmailAddresses.PurgeUniqueEmailAddress,
			)
		}
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

func (container *DIContainer) GetGRPCServer() *grpc.Server {
	if container.service.grpcServer == nil {
		container.service.grpcServer = grpc.NewServer()
		customergrpcproto.RegisterCustomerServer(container.service.grpcServer, container.getGRPCCustomerServer())
		reflection.Register(container.service.grpcServer)
	}

	return container.service.grpcServer
}
