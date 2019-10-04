package main

import (
	"database/sql"
	customercli "go-iddd/customer/api/cli"
	"go-iddd/service"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/urfave/cli"
)

var (
	logger         *service.Logger
	postgresDBConn *sql.DB
	diContainer    *service.DIContainer
)

func main() {
	bootstrap()
	mustRunCLIApp()
}

func bootstrap() {
	buildLogger()
	mustOpenPostgresDBConnection()
	mustRunPostgresDBMigrations()
	mustBuildDIContainer()
}

func buildLogger() {
	if logger == nil {
		logger = service.NewStandardLogger()
	}
}

func mustOpenPostgresDBConnection() {
	var err error

	if postgresDBConn == nil {
		dsn := "postgresql://goiddd:password123@localhost:5432/goiddd_local?sslmode=disable"

		if postgresDBConn, err = sql.Open("postgres", dsn); err != nil {
			logger.Errorf("failed to open Postgres DB handle: %s", err)
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

	migrator, err := migrate.NewWithDatabaseInstance("file://./service/dbmigrations", "postgres", driver)
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
		if diContainer, err = service.NewDIContainer(postgresDBConn); err != nil {
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
