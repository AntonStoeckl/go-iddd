package cmd

import (
	"database/sql"
	"go-iddd/service/customer/domain/customer/events"
	"go-iddd/service/lib/eventstore/postgres/database"
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

	migrator, err := database.NewMigrator(db, config.Postgres.MigrationsPath)
	if err != nil {
		return nil, err
	}

	err = migrator.Up()
	if err != nil {
		return nil, err
	}

	diContainer, err := NewDIContainer(db, events.UnmarshalCustomerEvent)

	if err != nil {
		return nil, err
	}

	return diContainer, nil
}
