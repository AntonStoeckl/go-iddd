package cmd

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/adapter/postgres/database"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/infrastructure/serialization"
	"github.com/AntonStoeckl/go-iddd/service/shared"
)

func Bootstrap(config *Config, logger *shared.Logger) (*DIContainer, error) {
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

	logger.Info("bootstrap: running DB migrations for customer ...")

	migratorCustomer, err := database.NewMigrator(db, config.Postgres.MigrationsPathCustomer)
	if err != nil {
		logger.Errorf("bootstrap: failed to create DB migrator for customer: %s", err)

		return nil, err
	}

	migratorCustomer.WithLogger(logger)

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
		customer.BuildUniqueEmailAddressAssertions,
	)
	if err != nil {
		logger.Errorf("bootstrap: failed to build the DI container: %s", err)

		return nil, err
	}

	return diContainer, nil
}
