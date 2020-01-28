package eventstore_test

import (
	"go-iddd/service/lib"
	"go-iddd/service/lib/infrastructure/eventstore/test"
	"math"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_PostgresEventStoreV2_LoadEventStream(t *testing.T) {
	Convey("Setup", t, func() {
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		eventStore := diContainer.GetPostgresEventStore()

		id := test.SomeID{Value: uuid.New().String()}
		streamID := lib.NewStreamID("customer" + "-" + id.ID())

		Convey("Given an empty event stream", func() {
			Convey("When it it loaded", func() {

				stream, err := eventStore.LoadEventStream(streamID, 0, math.MaxUint32)

				Convey("It should contain 0 events", func() {
					So(err, ShouldBeNil)
					So(stream, ShouldHaveLength, 0)
				})
			})
		})

		Convey("Given an event stream with 3 events", func() {
			event1 := test.CreateSomeEvent(id, 1)
			event2 := test.CreateSomeEvent(id, 2)
			event3 := test.CreateSomeEvent(id, 3)

			tx, err := db.Begin()
			So(err, ShouldBeNil)

			err = eventStore.AppendEventsToStream(
				streamID,
				lib.DomainEvents{event1, event2, event3},
				tx,
			)
			So(err, ShouldBeNil)

			errTx := tx.Commit()
			So(errTx, ShouldBeNil)

			Convey("When it it loaded", func() {
				stream, err := eventStore.LoadEventStream(streamID, 0, math.MaxUint32)

				Convey("It should contain 3 events", func() {
					So(err, ShouldBeNil)
					So(stream, ShouldHaveLength, 3)
				})
			})
		})
	})
}

func givenAnEventStream() {
	event1 := test.CreateSomeEvent(id, 1)
	event2 := test.CreateSomeEvent(id, 2)
	event3 := test.CreateSomeEvent(id, 3)

	tx, err := db.Begin()
	So(err, ShouldBeNil)

	err = eventStore.AppendEventsToStream(
		streamID,
		lib.DomainEvents{event1, event2, event3},
		tx,
	)
	So(err, ShouldBeNil)

	errTx := tx.Commit()
	So(errTx, ShouldBeNil)
}

func Test_PostgresEventStoreV2_PurgeEventStream(t *testing.T) {
	Convey("Setup", t, func() {
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		eventStore := diContainer.GetPostgresEventStore()

		Convey("Given an event stream with 3 events", func() {
			id := test.SomeID{Value: uuid.New().String()}
			streamID := lib.NewStreamID("customer" + "-" + id.ID())

			event1 := test.CreateSomeEvent(id, 1)
			event2 := test.CreateSomeEvent(id, 2)
			event3 := test.CreateSomeEvent(id, 3)

			tx, err := db.Begin()
			So(err, ShouldBeNil)

			err = eventStore.AppendEventsToStream(
				streamID,
				lib.DomainEvents{event1, event2, event3},
				tx,
			)
			So(err, ShouldBeNil)

			errTx := tx.Commit()
			So(errTx, ShouldBeNil)

			Convey("When the event stream is purged", func() {
				err := eventStore.PurgeEventStream(streamID)
				So(err, ShouldBeNil)

				Convey("It should contain 0 events", func() {
					stream, err := eventStore.LoadEventStream(streamID, 0, math.MaxUint32)
					So(err, ShouldBeNil)
					So(stream, ShouldHaveLength, 0)
				})
			})

			Convey("In case the DB connection was closed", func() {
				err := db.Close()
				So(err, ShouldBeNil)

				Convey("When the event stream is purged", func() {
					err := eventStore.PurgeEventStream(streamID)

					Convey("It should fail", func() {
						So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
					})
				})
			})
		})
	})
}
