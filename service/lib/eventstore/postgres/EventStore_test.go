package postgres_test

import (
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"go-iddd/service/lib/eventstore/postgres"
	"go-iddd/service/lib/eventstore/postgres/test"
	"math"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_EventStore_AppendEventsToStream(t *testing.T) {
	Convey("Setup", t, func() {
		diContainer, err := test.SetUpDIContainer()
		So(err, ShouldBeNil)
		db := diContainer.GetPostgresDBConn()
		eventStore := diContainer.GetEventStore()

		id := test.SomeID{Value: uuid.New().String()}
		streamID := es.NewStreamID("customer" + "-" + id.ID())

		Convey("Given an empty event stream", func() {
			Convey("When 2 events are appended", func() {
				appendedEvents := es.DomainEvents{
					test.CreateSomeEvent(id, 1),
					test.CreateSomeEvent(id, 2),
				}

				err := eventStore.AppendEventsToStream(streamID, appendedEvents)
				So(err, ShouldBeNil)

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

						err := eventStore.AppendEventsToStream(streamID, appendedEvents[2:4])
						So(err, ShouldBeNil)

						Convey("It should contain the expected 4 events", func() {
							eventStream, err := eventStore.LoadEventStream(streamID, 0, math.MaxUint32)
							So(err, ShouldBeNil)
							So(eventStream, ShouldHaveLength, 4)
							So(eventStream, ShouldResemble, appendedEvents)
						})
					})
				})

				// Clean up test artifacts
				err = eventStore.PurgeEventStream(streamID)
				So(err, ShouldBeNil)
			})

			Convey("When 2 events with conflicting stream version are appended", func() {
				event := test.CreateSomeEvent(id, 1)

				err = eventStore.AppendEventsToStream(
					streamID,
					es.DomainEvents{event, event},
				)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrConcurrencyConflict), ShouldBeTrue)
				})
			})

			Convey("When events which can't be marshaled to json are appended", func() {
				event := test.CreateBrokenMarshalingEvent(id, 1)

				err = eventStore.AppendEventsToStream(
					streamID,
					es.DomainEvents{event},
				)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrMarshalingFailed), ShouldBeTrue)
				})
			})
		})

		Convey("Given the DB connection was closed", func() {
			err := db.Close()
			So(err, ShouldBeNil)

			Convey("When events are appended", func() {
				event := test.CreateSomeEvent(id, 1)

				err = eventStore.AppendEventsToStream(
					streamID,
					es.DomainEvents{event},
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
			eventStore := postgres.NewEventStore(db, "unknown_table", test.UnmarshalMockEvents)

			id := test.SomeID{Value: uuid.New().String()}
			streamID := es.NewStreamID("customer" + "-" + id.ID())

			event1 := test.CreateSomeEvent(id, 1)
			event2 := test.CreateSomeEvent(id, 2)

			Convey("When events are appended", func() {
				err = eventStore.AppendEventsToStream(
					streamID,
					es.DomainEvents{event1, event2},
				)

				Convey("It should fail", func() {
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})
		})
	})
}

func Test_EventStore_LoadEventStream(t *testing.T) {
	Convey("Setup", t, func() {
		diContainer, err := test.SetUpDIContainer()
		So(err, ShouldBeNil)
		db := diContainer.GetPostgresDBConn()
		eventStore := diContainer.GetEventStore()

		id := test.SomeID{Value: uuid.New().String()}
		streamID := es.NewStreamID("customer" + "-" + id.ID())

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
			expectedEvents := es.DomainEvents{
				test.CreateSomeEvent(id, 1),
				test.CreateSomeEvent(id, 2),
				test.CreateSomeEvent(id, 3),
				test.CreateSomeEvent(id, 4),
			}

			err := eventStore.AppendEventsToStream(streamID, expectedEvents)
			So(err, ShouldBeNil)

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

			// Clean up test artifacts
			err = eventStore.PurgeEventStream(streamID)
			So(err, ShouldBeNil)
		})

		Convey("Given 3 events were appended with wrong order of stream versions", func() {
			expectedEvents := es.DomainEvents{
				test.CreateSomeEvent(id, 1),
				test.CreateSomeEvent(id, 2),
				test.CreateSomeEvent(id, 3),
			}

			err := eventStore.AppendEventsToStream(
				streamID,
				es.DomainEvents{
					expectedEvents[2],
					expectedEvents[0],
					expectedEvents[1],
				},
			)
			So(err, ShouldBeNil)

			Convey("When the event stream is loaded", func() {
				eventStream, err := eventStore.LoadEventStream(streamID, 0, math.MaxUint32)

				Convey("It should contain the expected 3 events properly ordered by stream versions", func() {
					So(err, ShouldBeNil)
					So(eventStream, ShouldHaveLength, 3)
					So(eventStream, ShouldResemble, expectedEvents)
				})
			})

			// Clean up test artifacts
			err = eventStore.PurgeEventStream(streamID)
			So(err, ShouldBeNil)
		})

		Convey("Given the eventStore contains an event which can't be unmarshaled", func() {
			event := test.CreateBrokenUnmarshalingEvent(id, 1)
			err := eventStore.AppendEventsToStream(streamID, es.DomainEvents{event})
			So(err, ShouldBeNil)

			Convey("When the event stream is loaded", func() {
				_, err := eventStore.LoadEventStream(streamID, 0, math.MaxUint32)

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
		diContainer, err := test.SetUpDIContainer()
		So(err, ShouldBeNil)
		db := diContainer.GetPostgresDBConn()
		eventStore := diContainer.GetEventStore()

		id := test.SomeID{Value: uuid.New().String()}
		streamID := es.NewStreamID("customer" + "-" + id.ID())

		Convey("Given an event stream with 3 events", func() {
			err := eventStore.AppendEventsToStream(
				streamID,
				es.DomainEvents{
					test.CreateSomeEvent(id, 1),
					test.CreateSomeEvent(id, 2),
					test.CreateSomeEvent(id, 3),
				},
			)
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
