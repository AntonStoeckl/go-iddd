package customers_test

import (
	"database/sql"
	"go-iddd/customer/domain"
	"go-iddd/customer/ports/secondary/customers"
	"go-iddd/shared"
	"go-iddd/shared/infrastructure/persistance/eventstore"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestEventSourcedRepository_StartSession(t *testing.T) {
	Convey("Given a Respository", t, func() {
		db, err := sql.Open("postgres", "postgresql://goiddd:password123@localhost:5432/goiddd_test?sslmode=disable")
		So(err, ShouldBeNil)
		eventStore := eventstore.NewPostgresEventStore(db, "eventstore", domain.UnmarshalDomainEvent)
		identityMap := customers.NewIdentityMap()
		repository := customers.NewEventSourcedRepository(eventStore, domain.ReconstituteCustomerFrom, identityMap)

		Convey("When a Session is started", func() {
			session, err := repository.StartSession()

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(session, ShouldNotBeNil)
				So(session, ShouldHaveSameTypeAs, &customers.EventSourcedRepositorySession{})
			})
		})

		Convey("And given the DB connection was closed", func() {
			err := db.Close()
			So(err, ShouldBeNil)

			Convey("When a Session is started", func() {
				session, err := repository.StartSession()

				Convey("It should fail", func() {
					So(xerrors.Is(err, shared.ErrTechnical), ShouldBeTrue)
					So(session, ShouldBeNil)
				})
			})
		})
	})
}
