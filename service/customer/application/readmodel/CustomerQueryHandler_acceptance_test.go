package readmodel_test

import (
	"database/sql"
	"go-iddd/service/cmd"
	"go-iddd/service/customer/application/writemodel/domain/customer/events"
	"go-iddd/service/lib/eventstore/postgres/database"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestQueryHandlerScenarios(t *testing.T) {
	diContainer := setUpDiContainerForCustomerQueryHandlerScenarios()
	customerEventStore := diContainer.GetCustomerEventStore()
	commandHandler := diContainer.GetCustomerQueryHandler()

	_, _ = customerEventStore, commandHandler

	Convey("Prepare test artifacts", t, func() {

	})
}

func setUpDiContainerForCustomerQueryHandlerScenarios() *cmd.DIContainer {
	config, err := cmd.NewConfigFromEnv()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", config.Postgres.DSN)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	migrator, err := database.NewMigrator(db, config.Postgres.MigrationsPath)
	if err != nil {
		panic(err)
	}

	err = migrator.Up()
	if err != nil {
		panic(err)
	}

	diContainer, err := cmd.NewDIContainer(db, events.UnmarshalCustomerEvent)
	if err != nil {
		panic(err)
	}

	return diContainer
}
