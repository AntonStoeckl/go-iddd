// +build test

package test

import (
	"database/sql"
	"go-iddd/service/lib"
	"go-iddd/service/lib/infrastructure/eventstore"

	"github.com/cockroachdb/errors"
)

const (
	eventStoreTableName = "eventstore"
)

type DIContainer struct {
	postgresDBConn       *sql.DB
	unmarshalDomainEvent lib.UnmarshalDomainEvent
	postgresEventStore   *eventstore.PostgresEventStoreV2
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

	return container, nil
}

func (container DIContainer) GetPostgresDBConn() *sql.DB {
	return container.postgresDBConn
}

func (container DIContainer) GetPostgresEventStore() *eventstore.PostgresEventStoreV2 {
	if container.postgresEventStore == nil {
		container.postgresEventStore = eventstore.NewPostgresEventStoreV2(
			container.postgresDBConn,
			eventStoreTableName,
			container.unmarshalDomainEvent,
		)
	}

	return container.postgresEventStore
}
