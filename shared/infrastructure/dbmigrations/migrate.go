// +xbuild generator

//go:generate go run migrate.go

package main

import (
	"database/sql"
	"fmt"
	"go-iddd/shared/infrastructure"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
)

func main() {
	var err error
	var config *infrastructure.Config
	var postgresDBConn *sql.DB

	config, err = infrastructure.NewConfigFromEnv()
	if err != nil {
		logrus.Fatalf("failed to get config from env - exiting: %s", err)
	}

	logger := infrastructure.NewStandardLogger()

	logger.Info("opening Postgres DB connection ...")

	if postgresDBConn, err = sql.Open("postgres", config.Postgres.DSN); err != nil {
		logger.Errorf("failed to open Postgres DB connection: %s", err)
	}

	if err := postgresDBConn.Ping(); err != nil {
		logger.Errorf("failed to connect to Postgres DB: %s", err)
	}

	driver, err := postgres.WithInstance(postgresDBConn, &postgres.Config{})
	if err != nil {
		logger.Errorf("failed to run migrations for Postgres DB: %s", err)
	}

	sourceURL := fmt.Sprintf("file://%s", config.Postgres.MigrationsPath)
	migrator, err := migrate.NewWithDatabaseInstance(sourceURL, "postgres", driver)
	if err != nil {
		logger.Errorf("failed to run migrations for Postgres DB: %s", err)
	}

	migrator.Log = logger

	if err := migrator.Up(); err != nil {
		if err != migrate.ErrNoChange {
			logger.Errorf("failed to run migrations for Postgres DB: %s", err)
		}
	}
}
