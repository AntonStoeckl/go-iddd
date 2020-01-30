package forstoringcustomerevents_test

import (
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/customer/infrastructure"
	"go-iddd/service/customer/infrastructure/secondary/forstoringcustomerevents"
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"go-iddd/service/lib/eventstore/mocked"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func Test_EventsourcedCustomers_With_Technical_Errors_From_EventStore(t *testing.T) {
	Convey("Setup", t, func() {
		id := values.GenerateCustomerID()
		var recordedEvents es.DomainEvents
		tx, err := infrastructure.MockTx()
		So(err, ShouldBeNil)
		eventStore := new(mocked.EventStore)
		customers := forstoringcustomerevents.NewEventsourcedCustomers(eventStore)

		Convey("Given a technical error from the EventStore when EventStream is called", func() {
			eventStore.
				On(
					"LoadEventStream",
					mock.AnythingOfType("es.StreamID"),
					mock.AnythingOfType("uint"),
					mock.AnythingOfType("uint"),
				).
				Return(nil, lib.ErrTechnical)

			Convey("When a Customer's eventStream is retrieved", func() {
				_, err := customers.EventStream(id)

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
					mock.AnythingOfType("*sql.Tx"),
				).
				Return(lib.ErrTechnical)

			Convey("When a Customer is registered", func() {
				err := customers.Register(id, recordedEvents, tx)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})

			Convey("When changes of a Customer are persisted", func() {
				err := customers.Persist(id, recordedEvents, tx)

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
