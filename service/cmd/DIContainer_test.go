package cmd

import (
	"database/sql"
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/shared"
	"github.com/AntonStoeckl/go-iddd/service/shared/es"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewDIContainer(t *testing.T) {
	Convey("When a DIContainer is created with valid input", t, func() {
		db, err := sql.Open("postgres", "postgresql://test:test@localhost:15432/test?sslmode=disable")
		So(err, ShouldBeNil)

		marshalDomainEvent := func(event es.DomainEvent) ([]byte, error) {
			return nil, nil
		}

		unmarshalDomainEvent := func(name string, payload []byte, streamVersion uint) (es.DomainEvent, error) {
			return nil, nil
		}

		diContainer, err := NewDIContainer(
			db,
			marshalDomainEvent,
			unmarshalDomainEvent,
			customer.BuildUniqueEmailAddressAssertions,
		)

		Convey("Then it should succeed", func() {
			So(err, ShouldBeNil)

			Convey("And it should expose the postgres DB connection", func() {
				So(diContainer.GetPostgresDBConn(), ShouldResemble, db)
			})
		})
	})

	Convey("When a DIContainer is created with a nil postgres DB connection", t, func() {
		var db *sql.DB

		_, err := NewDIContainer(
			db,
			func(event es.DomainEvent) ([]byte, error) { return nil, nil },
			func(name string, payload []byte, streamVersion uint) (es.DomainEvent, error) { return nil, nil },
			func(recordedEvents ...es.DomainEvent) customer.UniqueEmailAddressAssertions { return nil },
		)

		Convey("Then it should fail", func() {
			So(err, ShouldBeError)
			So(errors.Is(err, shared.ErrTechnical), ShouldBeTrue)
		})
	})
}
