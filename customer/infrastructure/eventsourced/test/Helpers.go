// +build test

package test

import (
	"database/sql"
	"go-iddd/customer/domain/events"
	"go-iddd/service"
	"go-iddd/shared/infrastructure/eventstore"

	. "github.com/smartystreets/goconvey/convey"
)

func SetUpDIContainer() *service.DIContainer {
	config, err := service.NewConfigFromEnv()
	So(err, ShouldBeNil)

	db, err := sql.Open("postgres", config.Postgres.DSN)
	So(err, ShouldBeNil)

	err = db.Ping()
	So(err, ShouldBeNil)

	migrator, err := eventstore.NewMigrator(db, config.Postgres.MigrationsPath)
	So(err, ShouldBeNil)

	err = migrator.Up()
	So(err, ShouldBeNil)

	diContainer, err := service.NewDIContainer(
		db,
		events.UnmarshalDomainEvent,
	)
	So(err, ShouldBeNil)

	return diContainer
}

func BeginTx(db *sql.DB) *sql.Tx {
	tx, err := db.Begin()
	So(err, ShouldBeNil)

	return tx
}
