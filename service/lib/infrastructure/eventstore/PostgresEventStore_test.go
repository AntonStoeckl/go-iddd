package eventstore_test

import (
	"go-iddd/service/lib"
	"go-iddd/service/lib/infrastructure/eventstore"
	"go-iddd/service/lib/infrastructure/eventstore/test"
	"math"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPostgresEventStore_StartSession(t *testing.T) {
	Convey("Given an EventStore", t, func() {
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		store := diContainer.GetPostgresEventStore()

		Convey("When a session is started", func() {
			tx := test.BeginTx(db)
			session := store.StartSession(tx)

			Convey("It should succeed", func() {
				So(session, ShouldNotBeNil)
				So(session, ShouldHaveSameTypeAs, &eventstore.PostgresEventStoreSession{})
				So(session, ShouldImplement, (*lib.EventStore)(nil))
			})
		})
	})
}

func TestPostgresEventStore_PurgeEventStream(t *testing.T) {
	Convey("Given an EventStore", t, func() {
		id := test.SomeID{ID: uuid.New().String()}
		streamID := lib.NewStreamID("customer" + "-" + id.String())
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		store := diContainer.GetPostgresEventStore()

		Convey("And given an event stream with 5 events", func() {
			event1 := test.CreateSomeEvent(id, 1)
			event2 := test.CreateSomeEvent(id, 2)
			event3 := test.CreateSomeEvent(id, 3)
			event4 := test.CreateSomeEvent(id, 4)
			event5 := test.CreateSomeEvent(id, 5)

			tx := test.BeginTx(db)
			session := store.StartSession(tx)

			err := session.AppendEventsToStream(
				streamID,
				lib.DomainEvents{event1, event2, event3, event4, event5},
			)
			So(err, ShouldBeNil)

			errTx := tx.Commit()
			So(errTx, ShouldBeNil)

			Convey("When the event stream is purged", func() {
				err := store.PurgeEventStream(streamID)
				So(err, ShouldBeNil)

				Convey("It should contain 0 events", func() {
					tx, err := db.Begin()
					So(err, ShouldBeNil)
					session := store.StartSession(tx)
					stream, err := session.LoadEventStream(streamID, 0, math.MaxUint32)
					So(err, ShouldBeNil)
					So(stream, ShouldHaveLength, 0)
					errTx = tx.Commit()
					So(errTx, ShouldBeNil)
				})
			})
		})

		Convey("And given the DB connection was closed", func() {
			err := db.Close()
			So(err, ShouldBeNil)

			Convey("When the event stream is purged", func() {
				err := store.PurgeEventStream(streamID)

				Convey("It should fail", func() {
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})
		})
	})
}
