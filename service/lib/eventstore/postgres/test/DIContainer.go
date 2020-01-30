// +build test

package test

import (
	"database/sql"
	"go-iddd/service/lib/es"
	"go-iddd/service/lib/eventstore/postgres"
	"go-iddd/service/lib/eventstore/postgres/database"
)

const (
	eventStoreTableName = "eventstore"
)

type DIContainer struct {
	postgresDBConn       *sql.DB
	unmarshalDomainEvent es.UnmarshalDomainEvent
	eventStore           *postgres.EventStore
}

func SetUpDIContainer() (*DIContainer, error) {
	config, err := NewConfigFromEnv()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("postgres", config.Postgres.DSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	migrator, err := database.NewMigrator(db, config.Postgres.MigrationsPath)
	if err != nil {
		return nil, err
	}

	err = migrator.Up()
	if err != nil {
		return nil, err
	}

	diContainer := &DIContainer{
		postgresDBConn:       db,
		unmarshalDomainEvent: UnmarshalMockEvents,
	}

	return diContainer, nil
}

func (container DIContainer) GetPostgresDBConn() *sql.DB {
	return container.postgresDBConn
}

func (container DIContainer) GetEventStore() *postgres.EventStore {
	if container.eventStore == nil {
		container.eventStore = postgres.NewEventStore(
			container.postgresDBConn,
			eventStoreTableName,
			container.unmarshalDomainEvent,
		)
	}

	return container.eventStore
}
