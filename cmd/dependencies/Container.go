package dependencies

import (
	"database/sql"
	"go-iddd/customer/application"
	"go-iddd/customer/domain"
	"go-iddd/customer/infrastructure/customers"
	"go-iddd/shared"
	"go-iddd/shared/infrastructure/persistance/eventstore"

	"github.com/cockroachdb/errors"
)

type Container struct {
	postgresDBConn         *sql.DB
	postgresEventStore     *eventstore.PostgresEventStore
	customerIdentityMap    *customers.IdentityMap
	customerRepository     application.StartsRepositorySessions
	customerCommandHandler *application.CommandHandler
}

func NewContainer(postgresDBConn *sql.DB) (*Container, error) {
	if postgresDBConn == nil {
		return nil, errors.Mark(errors.New("newContainer: postgres DB connection must not be nil"), shared.ErrTechnical)
	}

	return &Container{postgresDBConn: postgresDBConn}, nil
}

func (container Container) GetPostgresEventStore() *eventstore.PostgresEventStore {
	if container.postgresEventStore == nil {
		container.postgresEventStore = eventstore.NewPostgresEventStore(container.postgresDBConn, "eventstore", domain.UnmarshalDomainEvent)
	}

	return container.postgresEventStore
}

func (container Container) GetCustomerIdentityMap() *customers.IdentityMap {
	if container.customerIdentityMap == nil {
		container.customerIdentityMap = customers.NewIdentityMap()
	}

	return container.customerIdentityMap
}

func (container Container) GetCustomerRepository() application.StartsRepositorySessions {
	if container.customerRepository == nil {
		container.customerRepository = customers.NewEventSourcedRepository(
			container.GetPostgresEventStore(),
			domain.ReconstituteCustomerFrom,
			container.GetCustomerIdentityMap(),
		)
	}

	return container.customerRepository
}

func (container Container) GetCustomerCommandHandler() *application.CommandHandler {
	if container.customerCommandHandler == nil {
		container.customerCommandHandler = application.NewCommandHandler(
			container.GetCustomerRepository(),
			container.postgresDBConn,
		)
	}

	return container.customerCommandHandler
}
