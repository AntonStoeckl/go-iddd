package customers_test

import (
	"database/sql"
	"go-iddd/customer/domain"
	"go-iddd/customer/infrastructure/customers"
	"go-iddd/service"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEventSourcedRepository_StartSession(t *testing.T) {
	Convey("Given a Respository", t, func() {
		diContainer := setUpForEventSourcedRepository()
		db := diContainer.GetPostgresDBConn()
		repo := diContainer.GetCustomerRepository()

		Convey("When a Session is started", func() {
			tx, err := db.Begin()
			So(err, ShouldBeNil)

			session := repo.StartSession(tx)

			Convey("It should succeed", func() {
				So(session, ShouldNotBeNil)
				So(session, ShouldHaveSameTypeAs, &customers.EventSourcedRepositorySession{})
			})
		})
	})
}

func setUpForEventSourcedRepository() *service.DIContainer {
	config, err := service.NewConfigFromEnv()
	So(err, ShouldBeNil)

	db, err := sql.Open("postgres", config.Postgres.DSN)
	So(err, ShouldBeNil)

	diContainer, err := service.NewDIContainer(
		db,
		domain.UnmarshalDomainEvent,
		domain.ReconstituteCustomerFrom,
	)
	So(err, ShouldBeNil)

	return diContainer
}
