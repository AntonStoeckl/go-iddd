package main

import (
	"database/sql"
	"fmt"
	customercli "go-iddd/customer/api/cli"
	"go-iddd/customer/domain"
	"go-iddd/service"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	config         *service.Config
	logger         *service.Logger
	postgresDBConn *sql.DB
	diContainer    *service.DIContainer
)

func main() {
	bootstrap()
	mustRunCLIApp()
}

func bootstrap() {
	mustBuildConfig()
	buildLogger()
	mustOpenPostgresDBConnection()
	mustRunPostgresDBMigrations()
	mustBuildDIContainer()
}

func mustBuildConfig() {
	if config == nil {
		var err error

		config, err = service.NewConfigFromEnv()
		if err != nil {
			logrus.Fatalf("failed to get config from env - exiting: %s", err)
		}
	}
}

func buildLogger() {
	if logger == nil {
		logger = service.NewStandardLogger()
	}
}

func mustOpenPostgresDBConnection() {
	var err error

	if postgresDBConn == nil {
		if postgresDBConn, err = sql.Open("postgres", config.Postgres.DSN); err != nil {
			logger.Errorf("failed to open Postgres DB connection: %s", err)
			shutdown()
		}

		if err := postgresDBConn.Ping(); err != nil {
			logger.Errorf("failed to connect to Postgres DB: %s", err)
			shutdown()
		}
	}
}

func mustRunPostgresDBMigrations() {
	driver, err := postgres.WithInstance(postgresDBConn, &postgres.Config{})
	if err != nil {
		logger.Errorf("failed to run migrations for Postgres DB: %s", err)
		shutdown()
	}

	sourceURL := fmt.Sprintf("file://%s", config.Postgres.MigrationsPath)
	migrator, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		logger.Errorf("failed to run migrations for Postgres DB: %s", err)
		shutdown()
	}

	migrator.Log = logger

	if err := migrator.Up(); err != nil {
		if err != migrate.ErrNoChange {
			logger.Errorf("failed to run migrations for Postgres DB: %s", err)
			shutdown()
		}
	}
}

func mustBuildDIContainer() {
	var err error

	if diContainer == nil {
		diContainer, err = service.NewDIContainer(
			postgresDBConn,
			domain.UnmarshalDomainEvent,
			domain.ReconstituteCustomerFrom,
		)

		if err != nil {
			logger.Errorf("failed to build the DI container: %s", err)
			shutdown()
		}
	}
}

func mustRunCLIApp() {
	app := cli.NewApp()
	customerApp := customercli.NewCustomerApp(diContainer.GetCustomerCommandHandler())
	app.Commands = customerApp.Commands()

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}

func shutdown() {
	logger.Info("stopping services ...")

	if postgresDBConn != nil {
		//logger.Info("closing Postgres DB connection ...")
		if err := postgresDBConn.Close(); err != nil {
			logger.Warnf("failed to close the Postgres DB connection: %s", err)
		}
	}

	logger.Info("all services stopped - exiting")

	os.Exit(0)
}
