package eventstore_test

import (
	"database/sql"
	"go-iddd/service/lib"
	"go-iddd/service/lib/infrastructure/eventstore"
	"go-iddd/service/lib/infrastructure/eventstore/test"
	"math"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_PostgresEventStoreV2_AppendEventsToStream(t *testing.T) {
	Convey("Setup", t, func() {
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		eventStore := diContainer.GetPostgresEventStore()

		id := test.SomeID{Value: uuid.New().String()}
		streamID := lib.NewStreamID("customer" + "-" + id.ID())

		Convey("Given an empty event stream", func() {
			Convey("When 2 events are appended", func() {
				appendedEvents := lib.DomainEvents{
					test.CreateSomeEvent(id, 1),
					test.CreateSomeEvent(id, 2),
				}
				appendEventToStream(db, eventStore, streamID, appendedEvents[0])
				appendEventToStream(db, eventStore, streamID, appendedEvents[1])

				Convey("It should contain the expected 2 events", func() {
					eventStream, err := eventStore.LoadEventStream(streamID, 0, math.MaxUint32)
					So(err, ShouldBeNil)
					So(eventStream, ShouldHaveLength, 2)
					So(eventStream, ShouldResemble, appendedEvents)

					Convey("And when 2 additional events are appended", func() {
						appendedEvents = append(
							appendedEvents,
							test.CreateSomeEvent(id, 3),
							test.CreateSomeEvent(id, 4),
						)

						appendEventToStream(db, eventStore, streamID, appendedEvents[2])
						appendEventToStream(db, eventStore, streamID, appendedEvents[3])

						Convey("It should contain the expected 4 events", func() {
							eventStream, err := eventStore.LoadEventStream(streamID, 0, math.MaxUint32)
							So(err, ShouldBeNil)
							So(eventStream, ShouldHaveLength, 4)
							So(eventStream, ShouldResemble, appendedEvents)
						})
					})
				})
			})

			Convey("When 2 events with conflicting stream version are appended", func() {
				tx, err := db.Begin()
				So(err, ShouldBeNil)

				event := test.CreateSomeEvent(id, 1)

				err = eventStore.AppendEventsToStream(
					streamID,
					lib.DomainEvents{event, event},
					tx,
				)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrConcurrencyConflict), ShouldBeTrue)
				})
			})

			Convey("When boken events which can't be marshaled to json are appended", func() {
				tx, err := db.Begin()
				So(err, ShouldBeNil)

				event := test.CreateBrokenMarshalingEvent(id, 1)

				err = eventStore.AppendEventsToStream(
					streamID,
					lib.DomainEvents{event},
					tx,
				)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrMarshalingFailed), ShouldBeTrue)
				})
			})
		})

		Convey("Given the DB transaction was already committed", func() {
			tx, err := db.Begin()
			So(err, ShouldBeNil)

			err = tx.Commit()
			So(err, ShouldBeNil)

			Convey("When events are appended", func() {
				event := test.CreateSomeEvent(id, 1)

				err = eventStore.AppendEventsToStream(
					streamID,
					lib.DomainEvents{event},
					tx,
				)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})
		})
	})
}

func Test_PostgresEventStoreV2_LoadEventStream(t *testing.T) {
	Convey("Setup", t, func() {
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		eventStore := diContainer.GetPostgresEventStore()

		id := test.SomeID{Value: uuid.New().String()}
		streamID := lib.NewStreamID("customer" + "-" + id.ID())

		Convey("Given an empty event stream", func() {
			Convey("When it is loaded", func() {
				eventStream, err := eventStore.LoadEventStream(streamID, 0, math.MaxUint32)

				Convey("It should contain 0 events", func() {
					So(err, ShouldBeNil)
					So(eventStream, ShouldHaveLength, 0)
				})
			})
		})

		Convey("Given an event stream with 4 events", func() {
			expectedEvents := lib.DomainEvents{
				test.CreateSomeEvent(id, 1),
				test.CreateSomeEvent(id, 2),
				test.CreateSomeEvent(id, 3),
				test.CreateSomeEvent(id, 4),
			}
			appendEventToStream(db, eventStore, streamID, expectedEvents[0])
			appendEventToStream(db, eventStore, streamID, expectedEvents[1])
			appendEventToStream(db, eventStore, streamID, expectedEvents[2])
			appendEventToStream(db, eventStore, streamID, expectedEvents[3])

			Convey("When it is fully loaded", func() {
				eventStream, err := eventStore.LoadEventStream(streamID, 0, math.MaxUint32)

				Convey("It should contain the expected 4 events", func() {
					So(err, ShouldBeNil)
					So(eventStream, ShouldHaveLength, 4)
					So(eventStream, ShouldResemble, expectedEvents)
				})
			})

			Convey("When it is partially loaded (2 events starting from version 2)", func() {
				eventStream, err := eventStore.LoadEventStream(streamID, 2, 2)

				Convey("It should contain the expected 2 events", func() {
					So(err, ShouldBeNil)
					So(eventStream, ShouldHaveLength, 2)
					expectedEvents = expectedEvents[1:3]
					So(eventStream, ShouldResemble, expectedEvents)
				})
			})

		})

		Convey("Given 3 events were appended with wrong order of stream versions", func() {
			expectedEvents := lib.DomainEvents{
				test.CreateSomeEvent(id, 1),
				test.CreateSomeEvent(id, 2),
				test.CreateSomeEvent(id, 3),
			}
			appendEventToStream(db, eventStore, streamID, expectedEvents[2])
			appendEventToStream(db, eventStore, streamID, expectedEvents[0])
			appendEventToStream(db, eventStore, streamID, expectedEvents[1])

			Convey("When the event stream is loaded", func() {
				eventStream, err := eventStore.LoadEventStream(streamID, 0, math.MaxUint32)

				Convey("It should contain the expected 3 events properly ordered by stream versions", func() {
					So(err, ShouldBeNil)
					So(eventStream, ShouldHaveLength, 3)
					So(eventStream, ShouldResemble, expectedEvents)
				})
			})
		})

		Convey("In case the event store contains an event which can't be unmarshaled", func() {
			event := test.CreateBrokenUnmarshalingEvent(id, 1)
			appendEventToStream(db, eventStore, streamID, event)

			Convey("When the event stream is loaded", func() {
				_, err := eventStore.LoadEventStream(streamID, 0, math.MaxUint32)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrUnmarshalingFailed), ShouldBeTrue)
				})
			})
		})

		Convey("In case the DB connection was closed", func() {
			err := db.Close()
			So(err, ShouldBeNil)

			Convey("When the event stream is loaded", func() {
				_, err := eventStore.LoadEventStream(streamID, 0, math.MaxUint32)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})
		})
	})
}

func Test_PostgresEventStoreV2_PurgeEventStream(t *testing.T) {
	Convey("Setup", t, func() {
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		eventStore := diContainer.GetPostgresEventStore()

		id := test.SomeID{Value: uuid.New().String()}
		streamID := lib.NewStreamID("customer" + "-" + id.ID())

		Convey("Given an event stream with 3 events", func() {
			appendEventToStream(db, eventStore, streamID, test.CreateSomeEvent(id, 1))
			appendEventToStream(db, eventStore, streamID, test.CreateSomeEvent(id, 2))
			appendEventToStream(db, eventStore, streamID, test.CreateSomeEvent(id, 3))

			Convey("When the event stream is purged", func() {
				err := eventStore.PurgeEventStream(streamID)
				So(err, ShouldBeNil)

				Convey("It should contain 0 events", func() {
					stream, err := eventStore.LoadEventStream(streamID, 0, math.MaxUint32)
					So(err, ShouldBeNil)
					So(stream, ShouldHaveLength, 0)
				})
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
}

func appendEventToStream(
	db *sql.DB,
	eventStore *eventstore.PostgresEventStoreV2,
	streamID lib.StreamID,
	event lib.DomainEvent,
) {

	tx, err := db.Begin()
	So(err, ShouldBeNil)

	err = eventStore.AppendEventsToStream(
		streamID,
		lib.DomainEvents{event},
		tx,
	)
	So(err, ShouldBeNil)

	errTx := tx.Commit()
	So(errTx, ShouldBeNil)
}
