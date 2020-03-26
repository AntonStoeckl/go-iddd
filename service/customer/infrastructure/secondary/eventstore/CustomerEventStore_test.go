package eventstore_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/secondary/eventstore"
	customerMocked "github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/secondary/mocked"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	libMocked "github.com/AntonStoeckl/go-iddd/service/lib/eventstore/mocked"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func Test_CustomerEventStore_With_Technical_Errors_From_EventStore(t *testing.T) {
	Convey("Setup", t, func() {
		var recordedEvents es.DomainEvents

		id := values.GenerateCustomerID()

		eventStore := new(libMocked.EventStore)

		uniqEmailAddressChecker := new(customerMocked.ForCheckingUniqueEmailAddresses)

		dbMock, sqlMock, err := sqlmock.New()
		_ = sqlMock
		So(err, ShouldBeNil)

		customers := eventstore.NewCustomerEventStore(eventStore, uniqEmailAddressChecker, dbMock)

		Convey("Given a technical error from the EventStore when EventStreamFor is called", func() {
			eventStore.
				On(
					"LoadEventStream",
					mock.AnythingOfType("es.StreamID"),
					mock.AnythingOfType("uint"),
					mock.AnythingOfType("uint"),
				).
				Return(nil, lib.ErrTechnical)

			Convey("When a Customer's eventStream is retrieved", func() {
				_, err := customers.EventStreamFor(id)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})
		})

		Convey("Given a DB transaction can't be started", func() {
			sqlMock.ExpectBegin().WillReturnError(errors.Newf("mocked error: begin tx failed"))

			Convey("When a Customer is registered", func() {
				err := customers.CreateStreamFrom(recordedEvents, id)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})

			Convey("When changes of a Customer are persisted", func() {
				err := customers.Add(recordedEvents, id)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})

			So(sqlMock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("Given a DB transaction can't be committed", func() {
			eventStore.
				On(
					"AppendEventsToStream",
					mock.AnythingOfType("es.StreamID"),
					recordedEvents,
					mock.AnythingOfType("*sql.Tx"),
				).
				Return(nil)

			sqlMock.ExpectBegin()
			sqlMock.ExpectCommit().WillReturnError(errors.Newf("mocked error: commit tx failed"))

			Convey("When a Customer is registered", func() {
				err := customers.CreateStreamFrom(recordedEvents, id)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})

			Convey("When changes of a Customer are persisted", func() {
				err := customers.Add(recordedEvents, id)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})

			So(sqlMock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("Given a technical error from the EventStore when AppendEventsToStream is called", func() {
			eventStore.
				On(
					"AppendEventsToStream",
					mock.AnythingOfType("es.StreamID"),
					recordedEvents,
					mock.AnythingOfType("*sql.Tx"),
				).
				Return(lib.ErrTechnical)

			sqlMock.ExpectBegin()
			sqlMock.ExpectRollback()

			Convey("When a Customer is registered", func() {
				err := customers.CreateStreamFrom(recordedEvents, id)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})

			Convey("When changes of a Customer are persisted", func() {
				err := customers.Add(recordedEvents, id)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})

			So(sqlMock.ExpectationsWereMet(), ShouldBeNil)
		})

		Convey("Given a technical error from the EventStore when PurgeEventStream is called", func() {
			eventStore.
				On(
					"PurgeEventStream",
					mock.AnythingOfType("es.StreamID"),
				).
				Return(lib.ErrTechnical)

			Convey("When a Customer is deleted", func() {
				err := customers.Delete(id)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})
		})
	})
}
