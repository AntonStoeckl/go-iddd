package cmd

import (
	"database/sql"
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/shared"
	"github.com/AntonStoeckl/go-iddd/service/shared/es"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewDIContainer(t *testing.T) {
	marshalDomainEvent := func(event es.DomainEvent) ([]byte, error) {
		return nil, nil
	}

	unmarshalDomainEvent := func(name string, payload []byte, streamVersion uint) (es.DomainEvent, error) {
		return nil, nil
	}

	buildUniqueEmailAddressAssertions := func(recordedEvents ...es.DomainEvent) customer.UniqueEmailAddressAssertions {
		return nil
	}

	Convey("When a DIContainer is created with valid input", t, func() {
		db, err := sql.Open("postgres", "postgresql://test:test@localhost:15432/test?sslmode=disable")
		So(err, ShouldBeNil)

		logger := shared.NewNilLogger()
		config := MustBuildConfigFromEnv(logger)

		callback := func() {
			_ = MustBuildDIContainer(
				config,
				logger,
				UsePostgresDBConn(db),
				WithMarshalCustomerEvents(marshalDomainEvent),
				WithUnmarshalCustomerEvents(unmarshalDomainEvent),
				WithBuildUniqueEmailAddressAssertions(buildUniqueEmailAddressAssertions),
			)
		}

		Convey("Then it should succeed", func() {
			So(callback, ShouldNotPanic)
		})
	})

	Convey("When a DIContainer is created with a nil postgres DB connection", t, func() {
		var db *sql.DB

		logger := shared.NewNilLogger()
		config := MustBuildConfigFromEnv(logger)

		callback := func() {
			_ = MustBuildDIContainer(
				config,
				logger,
				UsePostgresDBConn(db),
				WithMarshalCustomerEvents(marshalDomainEvent),
				WithUnmarshalCustomerEvents(unmarshalDomainEvent),
				WithBuildUniqueEmailAddressAssertions(buildUniqueEmailAddressAssertions),
			)
		}

		Convey("Then it should panic", func() {
			So(callback, ShouldPanic)
		})
	})
}
