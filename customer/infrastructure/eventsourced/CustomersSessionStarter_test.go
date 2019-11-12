package eventsourced_test

import (
	"go-iddd/customer/infrastructure/eventsourced"
	"go-iddd/customer/infrastructure/eventsourced/test"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCustomersSessionStarter_StartSession(t *testing.T) {
	Convey("Given a Respository", t, func() {
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		repo := diContainer.GetCustomerRepository()

		Convey("When a Session is started", func() {
			tx := test.BeginTx(db)

			session := repo.StartSession(tx)

			Convey("It should succeed", func() {
				So(session, ShouldNotBeNil)
				So(session, ShouldHaveSameTypeAs, (*eventsourced.Customers)(nil))
			})
		})
	})
}
