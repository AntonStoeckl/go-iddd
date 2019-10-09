// +build test

package test

import (
	"database/sql"
	"go-iddd/shared/infrastructure/eventstore"
	"go-iddd/shared/infrastructure/eventstore/test/mocks"

	. "github.com/smartystreets/goconvey/convey"
)

func SetUpDIContainer() *DIContainer {
	config, err := NewConfigFromEnv()
	So(err, ShouldBeNil)

	db, err := sql.Open("postgres", config.Postgres.DSN)
	So(err, ShouldBeNil)

	err = db.Ping()
	So(err, ShouldBeNil)

	migrator, err := eventstore.NewMigrator(db, config.Postgres.MigrationsPath)
	So(err, ShouldBeNil)

	err = migrator.Up()
	So(err, ShouldBeNil)

	diContainer, err := NewDIContainer(db, mocks.Unmarshal)
	So(err, ShouldBeNil)

	return diContainer
}

func BeginTx(db *sql.DB) *sql.Tx {
	tx, err := db.Begin()
	So(err, ShouldBeNil)

	return tx
}
