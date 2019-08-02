package eventstore_test

import (
	"database/sql"
	"go-iddd/shared"
	"go-iddd/shared/infrastructure/persistance/eventstore"
	"go-iddd/shared/infrastructure/persistance/eventstore/mocks"
	"math"
	"testing"

	"github.com/google/uuid"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestPostgresEventStore_StartSession(t *testing.T) {
	Convey("Given an EventStore", t, func() {
		db, err := sql.Open("postgres", "postgresql://goiddd:password123@localhost:5432/goiddd?sslmode=disable")
		So(err, ShouldBeNil)
		store := eventstore.NewPostgresEventStore(db, "eventstore", mocks.Unmarshal)

		Convey("When a session is started", func() {
			tx, err := db.Begin()
			So(err, ShouldBeNil)

			session := store.StartSession(tx)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(session, ShouldNotBeNil)
				So(session, ShouldHaveSameTypeAs, &eventstore.PostgresEventStoreSession{})
				So(session, ShouldImplement, (*shared.EventStore)(nil))
			})
		})
	})
}

func TestPostgresEventStore_PurgeEventStream(t *testing.T) {
	Convey("Given an EventStore", t, func() {
		db, err := sql.Open("postgres", "postgresql://goiddd:password123@localhost:5432/goiddd?sslmode=disable")
		So(err, ShouldBeNil)
		store := eventstore.NewPostgresEventStore(db, "eventstore", mocks.Unmarshal)
		id := &mocks.SomeID{ID: uuid.New().String()}
		streamID := shared.NewStreamID("customer" + "-" + id.String())

		Convey("And given an event stream with 5 events", func() {
			event1 := mocks.CreateSomeEvent(id, 1)
			event2 := mocks.CreateSomeEvent(id, 2)
			event3 := mocks.CreateSomeEvent(id, 3)
			event4 := mocks.CreateSomeEvent(id, 4)
			event5 := mocks.CreateSomeEvent(id, 5)

			tx, err := db.Begin()
			So(err, ShouldBeNil)

			session := store.StartSession(tx)

			err = session.AppendEventsToStream(
				streamID,
				shared.DomainEvents{event1, event2, event3, event4, event5},
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
					So(xerrors.Is(err, shared.ErrTechnical), ShouldBeTrue)
				})
			})
		})
	})
}
