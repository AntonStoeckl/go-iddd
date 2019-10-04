package service

import (
	"database/sql"
	"go-iddd/customer/application"
	"go-iddd/customer/domain"
	"go-iddd/customer/infrastructure/customers"
	"go-iddd/shared"
	"go-iddd/shared/infrastructure/eventstore"

	"github.com/cockroachdb/errors"
)

const (
	eventStoreTableName = "eventstore"
)

type DIContainer struct {
	postgresDBConn         *sql.DB
	unmarshalDomainEvent   shared.UnmarshalDomainEvent
	customerFactory        func(eventStream shared.DomainEvents) (domain.Customer, error)
	postgresEventStore     *eventstore.PostgresEventStore
	customerIdentityMap    *customers.IdentityMap
	customerRepository     application.StartsRepositorySessions
	customerCommandHandler *application.CommandHandler
}

func NewDIContainer(
	postgresDBConn *sql.DB,
	unmarshalDomainEvent shared.UnmarshalDomainEvent,
	customerFactory func(eventStream shared.DomainEvents) (domain.Customer, error),
) (*DIContainer, error) {

	if postgresDBConn == nil {
		return nil, errors.Mark(errors.New("newContainer: postgres DB connection must not be nil"), shared.ErrTechnical)
	}

	container := &DIContainer{
		postgresDBConn:       postgresDBConn,
		unmarshalDomainEvent: unmarshalDomainEvent,
		customerFactory:      customerFactory,
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

func (container DIContainer) GetCustomerIdentityMap() *customers.IdentityMap {
	if container.customerIdentityMap == nil {
		container.customerIdentityMap = customers.NewIdentityMap()
	}

	return container.customerIdentityMap
}

func (container DIContainer) GetCustomerRepository() application.StartsRepositorySessions {
	if container.customerRepository == nil {
		container.customerRepository = customers.NewEventSourcedRepository(
			container.GetPostgresEventStore(),
			container.customerFactory,
			container.GetCustomerIdentityMap(),
		)
	}

	return container.customerRepository
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
