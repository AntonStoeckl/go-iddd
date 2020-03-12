package eventstore_test

import (
	"go-iddd/service/customer/application/readmodel/domain/customer"
	"go-iddd/service/customer/infrastructure/secondary/forreadingcustomereventstreams/eventstore"
	"go-iddd/service/lib"
	"go-iddd/service/lib/eventstore/mocked"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func Test_CustomerEventStore_With_Technical_Errors_From_EventStore(t *testing.T) {
	Convey("Setup", t, func() {
		id := customer.GenerateID()
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
	})
}
