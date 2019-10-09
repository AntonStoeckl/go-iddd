// +build test

package test

import (
	"database/sql"
	"go-iddd/shared/infrastructure/eventstore"
	"go-iddd/shared/infrastructure/eventstore/test/mocks"
)

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

	diContainer, err := NewDIContainer(db, mocks.Unmarshal)
	if err != nil {
		return nil, err
	}

	migrator, err := eventstore.NewMigrator(db, config.Postgres.MigrationsPath)
	if err != nil {
		return nil, err
	}

	err = migrator.Up()
	if err != nil {
		return nil, err
	}

	return diContainer, nil
}
