// +build test

package infrastructure

import (
	"database/sql"
	"go-iddd/service/cmd"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/lib/infrastructure/eventstore"

	. "github.com/smartystreets/goconvey/convey"
)

func SetUpDIContainer() *cmd.DIContainer {
	config, err := cmd.NewConfigFromEnv()
	So(err, ShouldBeNil)

	db, err := sql.Open("postgres", config.Postgres.DSN)
	So(err, ShouldBeNil)

	err = db.Ping()
	So(err, ShouldBeNil)

	migrator, err := eventstore.NewMigrator(db, config.Postgres.MigrationsPath)
	So(err, ShouldBeNil)

	err = migrator.Up()
	So(err, ShouldBeNil)

	diContainer, err := cmd.NewDIContainer(
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
