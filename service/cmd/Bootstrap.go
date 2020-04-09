package cmd

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/adapter/secondary/postgres"
	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/serialization"
	"github.com/AntonStoeckl/go-iddd/service/lib/eventstore/postgres/database"
)

func Bootstrap() (*DIContainer, error) {
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

	migratorEventstore, err := database.NewMigrator(db, config.Postgres.MigrationsPathEventstore)
	if err != nil {
		return nil, err
	}

	err = migratorEventstore.Up()
	if err != nil {
		return nil, err
	}

	migratorCustomer, err := postgres.NewMigrator(db, config.Postgres.MigrationsPathCustomer)
	if err != nil {
		return nil, err
	}

	err = migratorCustomer.Up()
	if err != nil {
		return nil, err
	}

	diContainer, err := NewDIContainer(db, serialization.MarshalCustomerEvent, serialization.UnmarshalCustomerEvent)

	if err != nil {
		return nil, err
	}

	return diContainer, nil
}
