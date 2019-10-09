package eventstore_test

import (
	"database/sql"
	"go-iddd/shared"
	"go-iddd/shared/infrastructure/eventstore"
	"go-iddd/shared/infrastructure/eventstore/test"
	"math"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestPostgresEventStoreSession_LoadEventStream(t *testing.T) {
	Convey("Given an empty event stream", t, func() {
		id := &test.SomeID{ID: uuid.New().String()}
		streamID := shared.NewStreamID("customer" + "-" + id.String())
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		store := diContainer.GetPostgresEventStore()

		Convey("When the event stream is loaded", func() {
			tx := test.BeginTx(db)
			session := store.StartSession(tx)

			stream, err := session.LoadEventStream(streamID, 0, math.MaxUint32)

			Convey("It should contain 0 events", func() {
				So(err, ShouldBeNil)
				So(stream, ShouldHaveLength, 0)
			})

			errTx := tx.Commit()
			So(errTx, ShouldBeNil)
		})

		closeDBForPostgresEventStoreSession(db)
	})

	Convey("Given an event stream with 5 events", t, func() {
		id := &test.SomeID{ID: uuid.New().String()}
		streamID := shared.NewStreamID("customer" + "-" + id.String())
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		store := diContainer.GetPostgresEventStore()
		event1 := test.CreateSomeEvent(id, 1)
		event2 := test.CreateSomeEvent(id, 2)
		event3 := test.CreateSomeEvent(id, 3)
		event4 := test.CreateSomeEvent(id, 4)
		event5 := test.CreateSomeEvent(id, 5)
		tx := test.BeginTx(db)
		session := store.StartSession(tx)
		err := session.AppendEventsToStream(
			streamID,
			shared.DomainEvents{event1, event2, event3, event4, event5},
		)
		So(err, ShouldBeNil)
		errTx := tx.Commit()
		So(errTx, ShouldBeNil)

		Convey("When the full event stream is loaded", func() {
			tx := test.BeginTx(db)
			session = store.StartSession(tx)

			stream, err := session.LoadEventStream(streamID, 0, math.MaxUint32)

			Convey("It should consist of the expected 5 events", func() {
				So(err, ShouldBeNil)
				So(
					stream,
					ShouldResemble,
					shared.DomainEvents{event1, event2, event3, event4, event5},
				)
			})

			errTx = tx.Commit()
			So(errTx, ShouldBeNil)
		})

		Convey("When 3 events starting from event with version 2 are loaded", func() {
			tx := test.BeginTx(db)
			session = store.StartSession(tx)

			stream, err := session.LoadEventStream(streamID, 2, 3)

			Convey("It should consist of the expected events of version 2, 3 and 4", func() {
				So(err, ShouldBeNil)
				So(
					stream,
					ShouldResemble,
					shared.DomainEvents{event2, event3, event4},
				)
			})

			errTx = tx.Commit()
			So(errTx, ShouldBeNil)
		})

		cleanUpArtefactsForPostgresEventStoreSession(store, streamID)
		closeDBForPostgresEventStoreSession(db)
	})

	Convey("Given an event in the stream can not be unmarshaled", t, func() {
		id := &test.SomeID{ID: uuid.New().String()}
		streamID := shared.NewStreamID("customer" + "-" + id.String())
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		store := diContainer.GetPostgresEventStore()
		event1 := test.CreateSomeEvent(id, 1)
		event2 := test.CreateBrokenUnmarshalingEvent(id, 2)
		tx := test.BeginTx(db)
		session := store.StartSession(tx)
		err := session.AppendEventsToStream(
			streamID,
			shared.DomainEvents{event1, event2},
		)
		So(err, ShouldBeNil)
		errTx := tx.Commit()
		So(errTx, ShouldBeNil)

		Convey("When the event stream is loaded", func() {
			tx := test.BeginTx(db)
			session = store.StartSession(tx)

			stream, err := session.LoadEventStream(streamID, 0, math.MaxUint32)

			Convey("It should fail", func() {
				So(xerrors.Is(err, shared.ErrUnmarshalingFailed), ShouldBeTrue)
				So(stream, ShouldHaveLength, 0)
			})

			errTx = tx.Commit()
			So(errTx, ShouldBeNil)
		})

		cleanUpArtefactsForPostgresEventStoreSession(store, streamID)
		closeDBForPostgresEventStoreSession(db)
	})

	Convey("Given the DB connection is already closed", t, func() {
		id := &test.SomeID{ID: uuid.New().String()}
		streamID := shared.NewStreamID("customer" + "-" + id.String())
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		store := diContainer.GetPostgresEventStore()
		tx := test.BeginTx(db)
		session := store.StartSession(tx)

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
		id := &test.SomeID{ID: uuid.New().String()}
		streamID := shared.NewStreamID("customer" + "-" + id.String())
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		store := diContainer.GetPostgresEventStore()

		event1 := test.CreateSomeEvent(id, 1)
		event2 := test.CreateSomeEvent(id, 2)
		event3 := test.CreateSomeEvent(id, 3)
		event4 := test.CreateSomeEvent(id, 4)
		event5 := test.CreateSomeEvent(id, 5)
		event6 := test.CreateSomeEvent(id, 6)

		Convey("When 3 events are appended", func() {
			tx := test.BeginTx(db)
			session := store.StartSession(tx)

			err := session.AppendEventsToStream(
				streamID,
				shared.DomainEvents{event1, event2, event3},
			)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)

				errTx := tx.Commit()
				So(errTx, ShouldBeNil)

				Convey("And the stream should consist of the expected 3 events", func() {
					tx := test.BeginTx(db)
					session := store.StartSession(tx)

					stream, err := session.LoadEventStream(streamID, 0, math.MaxUint32)
					So(err, ShouldBeNil)
					So(
						stream,
						ShouldResemble,
						shared.DomainEvents{event1, event2, event3},
					)

					errTx = tx.Commit()
					So(errTx, ShouldBeNil)

					Convey("And when 3 further events are appended", func() {
						tx := test.BeginTx(db)
						session := store.StartSession(tx)

						err := session.AppendEventsToStream(
							streamID,
							shared.DomainEvents{event4, event5, event6},
						)

						So(err, ShouldBeNil)

						errTx = tx.Commit()
						So(errTx, ShouldBeNil)

						Convey("The stream should consist of the expected 6 events", func() {
							tx := test.BeginTx(db)
							session := store.StartSession(tx)

							stream, err := session.LoadEventStream(streamID, 0, math.MaxUint32)
							So(err, ShouldBeNil)
							So(
								stream,
								ShouldResemble,
								shared.DomainEvents{event1, event2, event3, event4, event5, event6},
							)

							errTx = tx.Commit()
							So(errTx, ShouldBeNil)
						})
					})

					Convey("And when 3 further events are appended with a concurrency conflict", func() {
						tx := test.BeginTx(db)
						session := store.StartSession(tx)

						err := session.AppendEventsToStream(
							streamID,
							shared.DomainEvents{event4, event5, event6, event3}, // conflicting event last
						)

						Convey("It should fail", func() {
							So(xerrors.Is(err, shared.ErrConcurrencyConflict), ShouldBeTrue)

							errTx = tx.Rollback()
							So(errTx, ShouldBeNil)
						})
					})
				})
			})
		})

		cleanUpArtefactsForPostgresEventStoreSession(store, streamID)
		closeDBForPostgresEventStoreSession(db)

	})

	Convey("Given an event which can not be marshaled", t, func() {
		id := &test.SomeID{ID: uuid.New().String()}
		streamID := shared.NewStreamID("customer" + "-" + id.String())
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		store := diContainer.GetPostgresEventStore()

		event1 := test.CreateSomeEvent(id, 1)
		event2 := test.CreateBrokenMarshalingEvent(id, 2)

		Convey("When it is appended", func() {
			tx := test.BeginTx(db)
			session := store.StartSession(tx)

			err := session.AppendEventsToStream(
				streamID,
				shared.DomainEvents{event1, event2},
			)

			Convey("It should fail", func() {
				So(xerrors.Is(err, shared.ErrMarshalingFailed), ShouldBeTrue)
			})

			errTx := tx.Rollback()
			So(errTx, ShouldBeNil)
		})

		closeDBForPostgresEventStoreSession(db)
	})

	Convey("Given the session was already committed", t, func() {
		id := &test.SomeID{ID: uuid.New().String()}
		streamID := shared.NewStreamID("customer" + "-" + id.String())
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		store := diContainer.GetPostgresEventStore()
		tx := test.BeginTx(db)
		session := store.StartSession(tx)
		errTx := tx.Commit()
		So(errTx, ShouldBeNil)

		event1 := test.CreateSomeEvent(id, 1)
		event2 := test.CreateSomeEvent(id, 2)

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
		id := &test.SomeID{ID: uuid.New().String()}
		streamID := shared.NewStreamID("customer" + "-" + id.String())
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		store := eventstore.NewPostgresEventStore(db, "unknown_table", test.Unmarshal)

		event1 := test.CreateSomeEvent(id, 1)
		event2 := test.CreateSomeEvent(id, 2)

		Convey("When events are appended", func() {
			tx := test.BeginTx(db)
			session := store.StartSession(tx)

			err := session.AppendEventsToStream(
				streamID,
				shared.DomainEvents{event1, event2},
			)

			Convey("It should fail", func() {
				So(xerrors.Is(err, shared.ErrTechnical), ShouldBeTrue)
			})

			errTx := tx.Rollback()
			So(errTx, ShouldBeNil)
		})

		closeDBForPostgresEventStoreSession(db)
	})
}

/*** Test Helper Methods ***/

func cleanUpArtefactsForPostgresEventStoreSession(store *eventstore.PostgresEventStore, streamID *shared.StreamID) {
	err := store.PurgeEventStream(streamID)
	So(err, ShouldBeNil)
}

func closeDBForPostgresEventStoreSession(db *sql.DB) {
	err := db.Close()
	So(err, ShouldBeNil)
}
