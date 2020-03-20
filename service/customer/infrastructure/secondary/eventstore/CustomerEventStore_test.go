package eventstore_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/secondary/eventstore"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/AntonStoeckl/go-iddd/service/lib/eventstore/mocked"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func Test_CustomerEventStore_With_Technical_Errors_From_EventStore(t *testing.T) {
	Convey("Setup", t, func() {
		id := values.GenerateCustomerID()
		var recordedEvents es.DomainEvents
		eventStore := new(mocked.EventStore)
		customers := eventstore.NewCustomerEventStore(eventStore)

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

		Convey("Given a technical error from the EventStore when AppendEventsToStream is called", func() {
			eventStore.
				On(
					"AppendEventsToStream",
					mock.AnythingOfType("es.StreamID"),
					recordedEvents,
				).
				Return(lib.ErrTechnical)

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
