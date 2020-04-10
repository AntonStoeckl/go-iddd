package cmd

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/adapter/secondary/postgres"
	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/serialization"
	"github.com/AntonStoeckl/go-iddd/service/lib/eventstore/postgres/database"
)

func Bootstrap(config *Config, logger *Logger) (*DIContainer, error) {
	logger.Info("bootstrap: opening Postgres DB connection ...")

	db, err := sql.Open("postgres", config.Postgres.DSN)
	if err != nil {
		logger.Errorf("bootstrap: failed to open Postgres DB connection: %s", err)

		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Errorf("bootstrap: failed to connect to Postgres DB: %s", err)

		return nil, err
	}

	/***/

	logger.Info("bootstrap: running DB migrations for eventstore ...")

	migratorEventstore, err := database.NewMigrator(db, config.Postgres.MigrationsPathEventstore)
	if err != nil {
		logger.Errorf("bootstrap: failed to create DB migrator for eventstore: %s", err)

		return nil, err
	}

	err = migratorEventstore.Up()
	if err != nil {
		logger.Errorf("bootstrap: failed to run DB migrations for eventstore: %s", err)

		return nil, err
	}

	/***/

	logger.Info("bootstrap: running DB migrations for customer ...")

	migratorCustomer, err := postgres.NewMigrator(db, config.Postgres.MigrationsPathCustomer)
	if err != nil {
		logger.Errorf("bootstrap: failed to create DB migrator for customer: %s", err)

		return nil, err
	}

	err = migratorCustomer.Up()
	if err != nil {
		logger.Errorf("bootstrap: failed to run DB migrations for customer: %s", err)

		return nil, err
	}

	/***/

	logger.Info("bootstrap: building DI container ...")

	diContainer, err := NewDIContainer(
		db,
		serialization.MarshalCustomerEvent,
		serialization.UnmarshalCustomerEvent,
	)
	if err != nil {
		logger.Errorf("bootstrap: failed to build the DI container: %s", err)

		return nil, err
	}

	return diContainer, nil
}
