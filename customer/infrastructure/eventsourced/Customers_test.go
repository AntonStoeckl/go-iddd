package eventsourced_test

import (
	"errors"
	"fmt"
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/customer/infrastructure/eventsourced"
	"go-iddd/customer/infrastructure/eventsourced/test"
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestCustomers_Register(t *testing.T) {
	Convey("Given a Repository", t, func() {
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		repo := diContainer.GetCustomerRepository()

		Convey("And given a new Customer", func() {
			id := values.GenerateID()
			recordedEvents := registerCustomerForCustomersTest(id)

			Convey("When the Customer is registered", func() {
				tx := test.BeginTx(db)
				session := repo.StartSession(tx)

				err := session.Register(id, recordedEvents)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)
					err = tx.Commit()
					So(err, ShouldBeNil)

					Convey("And when the same Customer is registered again", func() {
						customer := registerCustomerForCustomersTest(id)
						tx := test.BeginTx(db)
						session := repo.StartSession(tx)

						err = session.Register(id, customer)

						Convey("It should fail", func() {
							So(xerrors.Is(err, shared.ErrDuplicate), ShouldBeTrue)
						})
					})
				})
			})

			Convey("And given the session was already committed", func() {
				recordedEvents := registerCustomerForCustomersTest(id)
				tx := test.BeginTx(db)
				session := repo.StartSession(tx)
				err := tx.Commit()
				So(err, ShouldBeNil)

				Convey("When the Customer is registered", func() {
					err = session.Register(id, recordedEvents)

					Convey("It should fail", func() {
						So(xerrors.Is(err, shared.ErrTechnical), ShouldBeTrue)
					})
				})
			})

			cleanUpArtefactsForCustomers(id)
		})

		Convey("And given an existing Customer", func() {
			id := values.GenerateID()
			recordedEvents := registerCustomerForCustomersTest(id)
			tx := test.BeginTx(db)
			session := repo.StartSession(tx)
			err := session.Register(id, recordedEvents)
			So(err, ShouldBeNil)
			err = tx.Commit()
			So(err, ShouldBeNil)

			Convey("When the same Customer is registered again", func() {
				recordedEvents := registerCustomerForCustomersTest(id)
				tx := test.BeginTx(db)
				session := repo.StartSession(tx)

				err = session.Register(id, recordedEvents)

				Convey("It should fail", func() {
					So(xerrors.Is(err, shared.ErrDuplicate), ShouldBeTrue)
				})
			})

			cleanUpArtefactsForCustomers(id)
		})
	})
}

func TestCustomers_Of(t *testing.T) {
	Convey("Given an existing Customer", t, func() {
		id := values.GenerateID()
		customer := registerCustomerForCustomersTest(id)
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		repo := diContainer.GetCustomerRepository()
		store := diContainer.GetPostgresEventStore()
		tx := test.BeginTx(db)
		session := repo.StartSession(tx)
		err := session.Register(id, customer)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)

		Convey("When the Customer is retrieved", func() {
			session := repo.StartSession(tx)

			customer, err := session.Of(id)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				So(customer.ID().Equals(id), ShouldBeTrue)
				So(customer.StreamVersion(), ShouldEqual, uint(1))
			})

			Convey("And when the customer can not be reconstituted", func() {
				expectedErr := errors.New("mocked error")
				customerFactory := func(eventStream shared.DomainEvents) (*domain.Customer, error) {
					return nil, expectedErr
				}

				repo := eventsourced.NewCustomersSessionStarter(store, customerFactory)
				tx := test.BeginTx(db)

				session := repo.StartSession(tx)

				customer, err := session.Of(id)

				Convey("It should fail", func() {
					So(xerrors.Is(err, expectedErr), ShouldBeTrue)
					So(customer, ShouldBeNil)
				})

				err = tx.Rollback()
				So(err, ShouldBeNil)
			})

		})

		Convey("And given the DB connection was closed", func() {
			tx := test.BeginTx(db)
			session := repo.StartSession(tx)

			err = db.Close()
			So(err, ShouldBeNil)

			Convey("When the Customer is retrieved", func() {
				customer, err := session.Of(id)

				Convey("It should fail", func() {
					So(xerrors.Is(err, shared.ErrTechnical), ShouldBeTrue)
					So(customer, ShouldBeNil)
				})
			})
		})

		cleanUpArtefactsForCustomers(id)
	})

	Convey("Given a not existing Customer", t, func() {
		id := values.GenerateID()
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		repo := diContainer.GetCustomerRepository()

		Convey("When the Customer is retrieved", func() {
			tx := test.BeginTx(db)
			session := repo.StartSession(tx)

			customer, err := session.Of(id)

			Convey("It should fail", func() {
				So(xerrors.Is(err, shared.ErrNotFound), ShouldBeTrue)
				So(customer, ShouldBeNil)
			})
		})
	})
}

func TestCustomers_Persist(t *testing.T) {
	Convey("Given a changed Customer", t, func() {
		id := values.GenerateID()
		recordedEvents := registerCustomerForCustomersTest(id)
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		repo := diContainer.GetCustomerRepository()
		tx := test.BeginTx(db)
		session := repo.StartSession(tx)
		err := session.Register(id, recordedEvents)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
		changeEmailAddress, err := commands.NewChangeEmailAddress(
			id.String(),
			fmt.Sprintf("john+%s+changed@doe.com", id.String()),
		)
		So(err, ShouldBeNil)
		customer, err := domain.ReconstituteCustomerFrom(recordedEvents)
		So(err, ShouldBeNil)

		recordedEvents = customer.ChangeEmailAddress(changeEmailAddress)

		Convey("When the Customer is persisted", func() {
			tx := test.BeginTx(db)
			session := repo.StartSession(tx)

			err = session.Persist(id, recordedEvents)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)
				err = tx.Commit()
				So(err, ShouldBeNil)

				tx := test.BeginTx(db)
				session := repo.StartSession(tx)
				customer, err := session.Of(id)
				So(err, ShouldBeNil)
				So(customer.ID().Equals(id), ShouldBeTrue)
				So(customer.StreamVersion(), ShouldEqual, uint(2))
				err = tx.Commit()
				So(err, ShouldBeNil)
			})
		})

		Convey("And given the session was already committed", func() {
			tx := test.BeginTx(db)
			session := repo.StartSession(tx)
			So(err, ShouldBeNil)

			err = tx.Commit()
			So(err, ShouldBeNil)

			Convey("When the Customer is persisted", func() {
				err = session.Persist(id, recordedEvents)

				Convey("It should fail", func() {
					So(xerrors.Is(err, shared.ErrTechnical), ShouldBeTrue)
				})
			})
		})

		cleanUpArtefactsForCustomers(id)
	})
}

/*** Test Helper Methods ***/

func registerCustomerForCustomersTest(id *values.ID) shared.DomainEvents {
	emailAddress := fmt.Sprintf("john+%s@doe.com", id.String())
	givenName := "John"
	familyName := "Doe"
	register, err := commands.NewRegister(id.String(), emailAddress, givenName, familyName)
	So(err, ShouldBeNil)

	recordedEvents := domain.RegisterCustomer(register)
	So(recordedEvents, ShouldHaveLength, 1)
	So(recordedEvents[0], ShouldHaveSameTypeAs, (*events.Registered)(nil))

	return recordedEvents
}

func cleanUpArtefactsForCustomers(id *values.ID) {
	diContainer := test.SetUpDIContainer()
	store := diContainer.GetPostgresEventStore()

	streamID := shared.NewStreamID("customer" + "-" + id.String())
	err := store.PurgeEventStream(streamID)
	So(err, ShouldBeNil)
}
