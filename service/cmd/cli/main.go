package main

import (
	"database/sql"
	"go-iddd/service/cmd"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/lib/eventstore/postgres/database"
	"os"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	config         *cmd.Config
	logger         *cmd.Logger
	postgresDBConn *sql.DB
	diContainer    *cmd.DIContainer
)

func main() {
	bootstrap()
	mustRunCLIApp()
}

func bootstrap() {
	mustBuildConfig()
	buildLogger()
	mustOpenPostgresDBConnection()
	mustRunDBMigrations()
	mustBuildDIContainer()
}

func mustBuildConfig() {
	if config == nil {
		var err error

		config, err = cmd.NewConfigFromEnv()
		if err != nil {
			logrus.Fatalf("failed to get config from env - exiting: %s", err)
		}
	}
}

func buildLogger() {
	if logger == nil {
		logger = cmd.NewStandardLogger()
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

func mustRunDBMigrations() {
	migrator, err := database.NewMigrator(postgresDBConn, config.Postgres.MigrationsPath)
	if err != nil {
		logger.Errorf("failed to create DB migrator: %s", err)
		shutdown()
		return
	}

	err = migrator.WithLogger(logger).Up()
	if err != nil {
		logger.Errorf("failed to run DB migrator: %s", err)
		shutdown()
	}
}

func mustBuildDIContainer() {
	var err error

	if diContainer == nil {
		diContainer, err = cmd.NewDIContainer(
			postgresDBConn,
			events.UnmarshalCustomerEvent,
		)

		if err != nil {
			logger.Errorf("failed to build the DI container: %s", err)
			shutdown()
		}
	}
}

func mustRunCLIApp() {
	app := cli.NewApp()
	customerApp := diContainer.GetCustomerApp()
	app.Commands = customerApp.Commands()

	if err := app.Run(os.Args); err != nil {
		logger.Fatal(err)
	}
}

func shutdown() {
	logger.Info("stopping services ...")

	if postgresDBConn != nil {
		if err := postgresDBConn.Close(); err != nil {
			logger.Warnf("failed to close the Postgres DB connection: %s", err)
		}
	}

	logger.Info("all services stopped - exiting")

	os.Exit(0)
}
