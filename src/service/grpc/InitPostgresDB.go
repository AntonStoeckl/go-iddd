package grpc

import (
	"database/sql"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/infrastructure/adapter/postgres/database"
	"github.com/AntonStoeckl/go-iddd/src/shared"
)

func MustInitPostgresDB(config *Config, logger *shared.Logger) *sql.DB {
	var err error

	logger.Info().Msg("bootstrapPostgresDB: opening Postgres DB connection ...")

	postgresDBConn, err := sql.Open("postgres", config.Postgres.DSN)
	if err != nil {
		logger.Panic().Msgf("bootstrapPostgresDB: failed to open Postgres DB connection: %s", err)
	}

	err = postgresDBConn.Ping()
	if err != nil {
		logger.Panic().Msgf("bootstrapPostgresDB: failed to connect to Postgres DB: %s", err)
	}

	/***/

	logger.Info().Msg("bootstrapPostgresDB: running DB migrations for customer ...")

	migratorCustomer, err := database.NewMigrator(postgresDBConn, config.Postgres.MigrationsPathCustomer)
	if err != nil {
		logger.Panic().Msgf("bootstrapPostgresDB: failed to create DB migrator for customer: %s", err)
	}

	migratorCustomer.WithLogger(logger)

	err = migratorCustomer.Up()
	if err != nil {
		logger.Panic().Msgf("bootstrapPostgresDB: failed to run DB migrations for customer: %s", err)
	}

	return postgresDBConn
}
