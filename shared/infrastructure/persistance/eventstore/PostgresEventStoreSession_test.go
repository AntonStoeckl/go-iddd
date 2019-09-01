package eventstore_test

import (
	"database/sql"
	"go-iddd/shared"
	"go-iddd/shared/infrastructure/persistance/eventstore"
	"go-iddd/shared/infrastructure/persistance/eventstore/mocks"
	"math"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestPostgresEventStoreSession_LoadEventStream(t *testing.T) {
	Convey("Given an empty event stream", t, func() {
		store, db, _, streamID := setUpForPostgresEventStoreSession()

		Convey("When the event stream is loaded", func() {
			session, errTx := store.StartSession()
			So(errTx, ShouldBeNil)

			stream, err := session.LoadEventStream(streamID, 0, math.MaxUint32)

			Convey("It should contain 0 events", func() {
				So(err, ShouldBeNil)
				So(stream, ShouldHaveLength, 0)
			})

			errTx = session.Commit()
			So(errTx, ShouldBeNil)
		})

		closeDBForPostgresEventStoreSession(db)
	})

	Convey("Given an event stream with 5 events", t, func() {
		store, db, id, streamID := setUpForPostgresEventStoreSession()
		event1 := mocks.CreateSomeEvent(id, 1)
		event2 := mocks.CreateSomeEvent(id, 2)
		event3 := mocks.CreateSomeEvent(id, 3)
		event4 := mocks.CreateSomeEvent(id, 4)
		event5 := mocks.CreateSomeEvent(id, 5)
		session, errTx := store.StartSession()
		So(errTx, ShouldBeNil)
		err := session.AppendEventsToStream(
			streamID,
			shared.DomainEvents{event1, event2, event3, event4, event5},
		)
		So(err, ShouldBeNil)
		errTx = session.Commit()
		So(errTx, ShouldBeNil)

		Convey("When the full event stream is loaded", func() {
			session, errTx = store.StartSession()
			So(errTx, ShouldBeNil)

			stream, err := session.LoadEventStream(streamID, 0, math.MaxUint32)

			Convey("It should consist of the expected 5 events", func() {
				So(err, ShouldBeNil)
				So(
					stream,
					ShouldResemble,
					shared.DomainEvents{event1, event2, event3, event4, event5},
				)
			})

			errTx = session.Commit()
			So(errTx, ShouldBeNil)
		})

		Convey("When 3 events starting from event with version 2 are loaded", func() {
			session, errTx = store.StartSession()
			So(errTx, ShouldBeNil)

			stream, err := session.LoadEventStream(streamID, 2, 3)

			Convey("It should consist of the expected events of version 2, 3 and 4", func() {
				So(err, ShouldBeNil)
				So(
					stream,
					ShouldResemble,
					shared.DomainEvents{event2, event3, event4},
				)
			})

			errTx = session.Commit()
			So(errTx, ShouldBeNil)
		})

		cleanUpArtefactsForPostgresEventStoreSession(store, streamID)
		closeDBForPostgresEventStoreSession(db)
	})

	Convey("Given an event in the stream can not be unmarshaled", t, func() {
		store, db, id, streamID := setUpForPostgresEventStoreSession()
		event1 := mocks.CreateSomeEvent(id, 1)
		event2 := mocks.CreateBrokenUnmarshalingEvent(id, 2)
		session, errTx := store.StartSession()
		So(errTx, ShouldBeNil)
		err := session.AppendEventsToStream(
			streamID,
			shared.DomainEvents{event1, event2},
		)
		So(err, ShouldBeNil)
		errTx = session.Commit()
		So(errTx, ShouldBeNil)

		Convey("When the event stream is loaded", func() {
			session, errTx = store.StartSession()
			So(errTx, ShouldBeNil)

			stream, err := session.LoadEventStream(streamID, 0, math.MaxUint32)

			Convey("It should fail", func() {
				So(xerrors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
				So(stream, ShouldHaveLength, 0)
			})

			errTx = session.Commit()
			So(errTx, ShouldBeNil)
		})

		cleanUpArtefactsForPostgresEventStoreSession(store, streamID)
		closeDBForPostgresEventStoreSession(db)
	})

	Convey("Given the DB connection is already closed", t, func() {
		store, db, _, streamID := setUpForPostgresEventStoreSession()
		session, errTx := store.StartSession()
		So(errTx, ShouldBeNil)

		closeDBForPostgresEventStoreSession(db)

		Convey("When the event stream is loaded", func() {
			_, err := session.LoadEventStream(streamID, 0, math.MaxUint32)

			Convey("It should fail", func() {
				So(xerrors.Is(err, shared.ErrTechnical), ShouldBeTrue)
			})
		})
	})
}

func TestPostgresEventStoreSession_AppendEventsToStream(t *testing.T) {
	Convey("Given an empty event stream", t, func() {
		store, db, id, streamID := setUpForPostgresEventStoreSession()

		event1 := mocks.CreateSomeEvent(id, 1)
		event2 := mocks.CreateSomeEvent(id, 2)
		event3 := mocks.CreateSomeEvent(id, 3)
		event4 := mocks.CreateSomeEvent(id, 4)
		event5 := mocks.CreateSomeEvent(id, 5)
		event6 := mocks.CreateSomeEvent(id, 6)

		Convey("When 3 events are appended", func() {
			session, errTx := store.StartSession()
			So(errTx, ShouldBeNil)

			err := session.AppendEventsToStream(
				streamID,
				shared.DomainEvents{event1, event2, event3},
			)

			So(err, ShouldBeNil)

			errTx = session.Commit()
			So(errTx, ShouldBeNil)

			Convey("The stream should consist of the expected 3 events", func() {
				session, errTx := store.StartSession()
				So(errTx, ShouldBeNil)

				stream, err := session.LoadEventStream(streamID, 0, math.MaxUint32)
				So(err, ShouldBeNil)
				So(
					stream,
					ShouldResemble,
					shared.DomainEvents{event1, event2, event3},
				)

				errTx = session.Commit()
				So(errTx, ShouldBeNil)

				Convey("And when 3 further events are appended", func() {
					session, errTx := store.StartSession()
					So(errTx, ShouldBeNil)

					err := session.AppendEventsToStream(
						streamID,
						shared.DomainEvents{event4, event5, event6},
					)

					So(err, ShouldBeNil)

					errTx = session.Commit()
					So(errTx, ShouldBeNil)

					Convey("The stream should consist of the expected 6 events", func() {
						session, errTx := store.StartSession()
						So(errTx, ShouldBeNil)

						stream, err := session.LoadEventStream(streamID, 0, math.MaxUint32)
						So(err, ShouldBeNil)
						So(
							stream,
							ShouldResemble,
							shared.DomainEvents{event1, event2, event3, event4, event5, event6},
						)

						errTx = session.Commit()
						So(errTx, ShouldBeNil)
					})
				})

				Convey("And when 3 further events are appended with a concurrency conflict", func() {
					session, errTx := store.StartSession()
					So(errTx, ShouldBeNil)

					err := session.AppendEventsToStream(
						streamID,
						shared.DomainEvents{event4, event5, event6, event3}, // conflicting event last
					)

					Convey("It should fail", func() {
						So(xerrors.Is(err, shared.ErrConcurrencyConflict), ShouldBeTrue)

						errTx = session.Rollback()
						So(errTx, ShouldBeNil)
					})
				})
			})
		})

		cleanUpArtefactsForPostgresEventStoreSession(store, streamID)
		closeDBForPostgresEventStoreSession(db)

	})

	Convey("Given an event which can not be marshaled", t, func() {
		store, db, id, streamID := setUpForPostgresEventStoreSession()

		event1 := mocks.CreateSomeEvent(id, 1)
		event2 := mocks.CreateBrokenMarshalingEvent(id, 2)

		Convey("When it is appended", func() {
			session, errTx := store.StartSession()
			So(errTx, ShouldBeNil)

			err := session.AppendEventsToStream(
				streamID,
				shared.DomainEvents{event1, event2},
			)

			Convey("It should fail", func() {
				So(xerrors.Is(err, shared.ErrMarshalingFailed), ShouldBeTrue)
			})

			errTx = session.Rollback()
			So(errTx, ShouldBeNil)
		})

		closeDBForPostgresEventStoreSession(db)
	})

	Convey("Given the session was already committed", t, func() {
		store, db, id, streamID := setUpForPostgresEventStoreSession()
		session, errTx := store.StartSession()
		So(errTx, ShouldBeNil)
		errTx = session.Commit()
		So(errTx, ShouldBeNil)

		event1 := mocks.CreateSomeEvent(id, 1)
		event2 := mocks.CreateSomeEvent(id, 2)

		Convey("When events are appended", func() {
			err := session.AppendEventsToStream(
				streamID,
				shared.DomainEvents{event1, event2},
			)

			Convey("It should fail", func() {
				So(xerrors.Is(err, shared.ErrTechnical), ShouldBeTrue)
			})
		})

		cleanUpArtefactsForPostgresEventStoreSession(store, streamID)
		closeDBForPostgresEventStoreSession(db)
	})

	Convey("Given the DB table does not exist", t, func() {
		store, db, id, streamID := setUpForPostgresEventStoreSession()
		store = eventstore.NewPostgresEventStore(db, "unknown_table", mocks.Unmarshal)

		event1 := mocks.CreateSomeEvent(id, 1)
		event2 := mocks.CreateSomeEvent(id, 2)

		Convey("When events are appended", func() {
			session, errTx := store.StartSession()
			So(errTx, ShouldBeNil)

			err := session.AppendEventsToStream(
				streamID,
				shared.DomainEvents{event1, event2},
			)

			Convey("It should fail", func() {
				So(xerrors.Is(err, shared.ErrTechnical), ShouldBeTrue)
			})

			errTx = session.Rollback()
			So(errTx, ShouldBeNil)
		})

		closeDBForPostgresEventStoreSession(db)
	})
}

func TestPostgresEventStoreSession_Commit(t *testing.T) {
	Convey("Given a session which was already rolled back", t, func() {
		store, db, _, _ := setUpForPostgresEventStoreSession()
		session, err := store.StartSession()
		So(err, ShouldBeNil)

		err = session.Rollback()
		So(err, ShouldBeNil)

		Convey("When the session is committed", func() {
			err = session.Commit()

			Convey("It should fail", func() {
				So(xerrors.Is(err, shared.ErrTechnical), ShouldBeTrue)
			})
		})

		closeDBForPostgresEventStoreSession(db)
	})
}

func TestPostgresEventStoreSession_Rollback(t *testing.T) {
	Convey("Given a session that was already committed", t, func() {
		store, db, _, _ := setUpForPostgresEventStoreSession()
		session, err := store.StartSession()
		So(err, ShouldBeNil)

		err = session.Commit()
		So(err, ShouldBeNil)

		Convey("When the session is rolled back", func() {
			err = session.Rollback()

			Convey("It should fail", func() {
				So(xerrors.Is(err, shared.ErrTechnical), ShouldBeTrue)
			})
		})

		closeDBForPostgresEventStoreSession(db)
	})
}

/*** Test Helper Methods ***/

func setUpForPostgresEventStoreSession() (*eventstore.PostgresEventStore, *sql.DB, *mocks.SomeID, *shared.StreamID) {
	db, err := sql.Open("postgres", "postgresql://goiddd:password123@localhost:5432/goiddd_test?sslmode=disable")
	So(err, ShouldBeNil)
	es := eventstore.NewPostgresEventStore(db, "eventstore", mocks.Unmarshal)
	id := &mocks.SomeID{ID: uuid.New().String()}
	streamID := shared.NewStreamID("customer" + "-" + id.String())

	return es, db, id, streamID
}

func cleanUpArtefactsForPostgresEventStoreSession(store *eventstore.PostgresEventStore, streamID *shared.StreamID) {
	err := store.PurgeEventStream(streamID)
	So(err, ShouldBeNil)
}

func closeDBForPostgresEventStoreSession(db *sql.DB) {
	err := db.Close()
	So(err, ShouldBeNil)
}
