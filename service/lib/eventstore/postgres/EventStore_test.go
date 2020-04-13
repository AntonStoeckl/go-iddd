package postgres_test

import (
	"database/sql"
	"math"
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/AntonStoeckl/go-iddd/service/lib/eventstore/postgres"
	"github.com/AntonStoeckl/go-iddd/service/lib/eventstore/postgres/test"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_EventStore_AppendEventsToStream(t *testing.T) {
	Convey("Setup", t, func() {
		var err error
		var tx *sql.Tx
		var eventStream es.EventStream

		diContainer, err := test.SetUpDIContainer()
		So(err, ShouldBeNil)
		db := diContainer.GetPostgresDBConn()
		eventStore := diContainer.GetEventStore()

		id := test.SomeID{Value: uuid.New().String()}
		streamID := es.NewStreamID("customer" + "-" + id.ID())

		Convey("Given an empty event stream", func() {
			Convey("When 2 events are appended", func() {
				event1 := test.CreateSomeEvent(id, 1)
				event2 := test.CreateSomeEvent(id, 2)

				tx, err = db.Begin()
				So(err, ShouldBeNil)

				err = eventStore.AppendEventsToStream(streamID, es.RecordedEvents{event1, event2}, tx)
				So(err, ShouldBeNil)

				err = tx.Commit()
				So(err, ShouldBeNil)

				Convey("It should contain the expected 2 events", func() {
					eventStream, err = eventStore.LoadEventStream(streamID, 0, math.MaxUint32)
					So(err, ShouldBeNil)
					So(eventStream, ShouldHaveLength, 2)
					So(eventStream, ShouldResemble, es.EventStream{event1, event2})

					Convey("And when 2 additional events are appended", func() {
						event3 := test.CreateSomeEvent(id, 3)
						event4 := test.CreateSomeEvent(id, 4)

						tx, err = db.Begin()
						So(err, ShouldBeNil)

						err = eventStore.AppendEventsToStream(streamID, es.RecordedEvents{event3, event4}, tx)
						So(err, ShouldBeNil)

						err = tx.Commit()
						So(err, ShouldBeNil)

						Convey("It should contain the expected 4 events", func() {
							eventStream, err = eventStore.LoadEventStream(streamID, 0, math.MaxUint32)
							So(err, ShouldBeNil)
							So(eventStream, ShouldHaveLength, 4)
							So(eventStream, ShouldResemble, es.EventStream{event1, event2, event3, event4})
						})
					})
				})

				// Clean up test artifacts
				err = eventStore.PurgeEventStream(streamID)
				So(err, ShouldBeNil)
			})

			Convey("When 2 events with conflicting stream version are appended", func() {
				event := test.CreateSomeEvent(id, 1)

				tx, err = db.Begin()
				So(err, ShouldBeNil)

				err = eventStore.AppendEventsToStream(
					streamID,
					es.RecordedEvents{event, event},
					tx,
				)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrConcurrencyConflict), ShouldBeTrue)

					err = tx.Rollback()
					So(err, ShouldBeNil)
				})
			})

			Convey("When events which can't be marshaled to json are appended", func() {
				event := test.CreateBrokenMarshalingEvent(id, 1)

				tx, err = db.Begin()
				So(err, ShouldBeNil)

				err = eventStore.AppendEventsToStream(
					streamID,
					es.RecordedEvents{event},
					tx,
				)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrMarshalingFailed), ShouldBeTrue)

					err = tx.Rollback()
					So(err, ShouldBeNil)
				})
			})
		})

		Convey("Given the DB transaction was closed", func() {
			tx, err = db.Begin()
			So(err, ShouldBeNil)

			err = tx.Rollback()
			So(err, ShouldBeNil)

			Convey("When events are appended", func() {
				event := test.CreateSomeEvent(id, 1)

				err = eventStore.AppendEventsToStream(
					streamID,
					es.RecordedEvents{event},
					tx,
				)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})
		})

		Convey("Given the DB table does not exist", func() {
			diContainer, err := test.SetUpDIContainer()
			So(err, ShouldBeNil)
			db := diContainer.GetPostgresDBConn()
			eventStore := postgres.NewEventStore(db, "unknown_table", test.MarshalMockEvents, test.UnmarshalMockEvents)

			id := test.SomeID{Value: uuid.New().String()}
			streamID := es.NewStreamID("customer" + "-" + id.ID())

			event1 := test.CreateSomeEvent(id, 1)
			event2 := test.CreateSomeEvent(id, 2)

			Convey("When events are appended", func() {
				tx, err = db.Begin()
				So(err, ShouldBeNil)

				err = eventStore.AppendEventsToStream(
					streamID,
					es.RecordedEvents{event1, event2},
					tx,
				)

				Convey("It should fail", func() {
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)

					err = tx.Rollback()
					So(err, ShouldBeNil)
				})
			})
		})
	})
}

func Test_EventStore_LoadEventStream(t *testing.T) {
	Convey("Setup", t, func() {
		var err error
		var tx *sql.Tx
		var eventStream es.EventStream

		diContainer, err := test.SetUpDIContainer()
		So(err, ShouldBeNil)
		db := diContainer.GetPostgresDBConn()
		eventStore := diContainer.GetEventStore()

		id := test.SomeID{Value: uuid.New().String()}
		streamID := es.NewStreamID("customer" + "-" + id.ID())

		Convey("Given an empty event stream", func() {
			Convey("When it is loaded", func() {
				eventStream, err = eventStore.LoadEventStream(streamID, 0, math.MaxUint32)

				Convey("It should contain 0 events", func() {
					So(err, ShouldBeNil)
					So(eventStream, ShouldHaveLength, 0)
				})
			})
		})

		Convey("Given an event stream with 4 events", func() {
			event1 := test.CreateSomeEvent(id, 1)
			event2 := test.CreateSomeEvent(id, 2)
			event3 := test.CreateSomeEvent(id, 3)
			event4 := test.CreateSomeEvent(id, 4)

			tx, err = db.Begin()
			So(err, ShouldBeNil)

			err = eventStore.AppendEventsToStream(streamID, es.RecordedEvents{event1, event2, event3, event4}, tx)
			So(err, ShouldBeNil)

			err = tx.Commit()
			So(err, ShouldBeNil)

			Convey("When it is fully loaded", func() {
				eventStream, err = eventStore.LoadEventStream(streamID, 0, math.MaxUint32)

				Convey("It should contain the expected 4 events", func() {
					So(err, ShouldBeNil)
					So(eventStream, ShouldHaveLength, 4)
					So(eventStream, ShouldResemble, es.EventStream{event1, event2, event3, event4})
				})
			})

			Convey("When it is partially loaded (2 events starting from version 2)", func() {
				eventStream, err = eventStore.LoadEventStream(streamID, 2, 2)

				Convey("It should contain the expected 2 events", func() {
					So(err, ShouldBeNil)
					So(eventStream, ShouldHaveLength, 2)
					So(eventStream, ShouldResemble, es.EventStream{event2, event3})
				})
			})

			// Clean up test artifacts
			err = eventStore.PurgeEventStream(streamID)
			So(err, ShouldBeNil)
		})

		Convey("Given 3 events were appended with wrong order of stream versions", func() {
			event1 := test.CreateSomeEvent(id, 1)
			event2 := test.CreateSomeEvent(id, 2)
			event3 := test.CreateSomeEvent(id, 3)

			tx, err = db.Begin()
			So(err, ShouldBeNil)

			err = eventStore.AppendEventsToStream(streamID, es.RecordedEvents{event3, event1, event2}, tx)
			So(err, ShouldBeNil)

			err = tx.Commit()
			So(err, ShouldBeNil)

			Convey("When the event stream is loaded", func() {
				eventStream, err = eventStore.LoadEventStream(streamID, 0, math.MaxUint32)

				Convey("It should contain the expected 3 events properly ordered by stream versions", func() {
					So(err, ShouldBeNil)
					So(eventStream, ShouldHaveLength, 3)
					So(eventStream, ShouldResemble, es.EventStream{event1, event2, event3})
				})
			})

			// Clean up test artifacts
			err = eventStore.PurgeEventStream(streamID)
			So(err, ShouldBeNil)
		})

		Convey("Given the eventStore contains an event which can't be unmarshaled", func() {
			event := test.CreateBrokenUnmarshalingEvent(id, 1)

			tx, err = db.Begin()
			So(err, ShouldBeNil)

			err = eventStore.AppendEventsToStream(streamID, es.RecordedEvents{event}, tx)
			So(err, ShouldBeNil)

			err = tx.Commit()
			So(err, ShouldBeNil)

			Convey("When the event stream is loaded", func() {
				_, err = eventStore.LoadEventStream(streamID, 0, math.MaxUint32)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrUnmarshalingFailed), ShouldBeTrue)
				})
			})

			// Clean up test artifacts
			err = eventStore.PurgeEventStream(streamID)
			So(err, ShouldBeNil)
		})

		Convey("Given the DB connection was closed", func() {
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

func Test_EventStore_PurgeEventStream(t *testing.T) {
	Convey("Setup", t, func() {
		var err error
		var tx *sql.Tx

		diContainer, err := test.SetUpDIContainer()
		So(err, ShouldBeNil)
		db := diContainer.GetPostgresDBConn()
		eventStore := diContainer.GetEventStore()

		id := test.SomeID{Value: uuid.New().String()}
		streamID := es.NewStreamID("customer" + "-" + id.ID())

		Convey("Given an event stream with 3 events", func() {
			tx, err = db.Begin()
			So(err, ShouldBeNil)

			event1 := test.CreateSomeEvent(id, 1)
			event2 := test.CreateSomeEvent(id, 2)
			event3 := test.CreateSomeEvent(id, 3)

			err = eventStore.AppendEventsToStream(streamID, es.RecordedEvents{event1, event2, event3}, tx)
			So(err, ShouldBeNil)

			err = tx.Commit()
			So(err, ShouldBeNil)

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
