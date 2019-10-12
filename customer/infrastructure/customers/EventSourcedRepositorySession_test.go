package customers_test

import (
	"errors"
	"fmt"
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
	"go-iddd/customer/infrastructure/customers"
	"go-iddd/customer/infrastructure/customers/test"
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestEventSourcedRepositorySession_Register(t *testing.T) {
	Convey("Given a Repository", t, func() {
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		repo := diContainer.GetCustomerRepository()

		Convey("And given a new Customer", func() {
			id := values.GenerateID()
			customer := buildRegisteredCustomerWith(id)

			Convey("When the Customer is registered", func() {
				tx := test.BeginTx(db)
				session := repo.StartSession(tx)

				err := session.Register(customer)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)
					err = tx.Commit()
					So(err, ShouldBeNil)

					Convey("And when the same Customer is registered again", func() {
						customer := buildRegisteredCustomerWith(id)
						tx := test.BeginTx(db)
						session := repo.StartSession(tx)

						err = session.Register(customer)

						Convey("It should fail", func() {
							So(xerrors.Is(err, shared.ErrDuplicate), ShouldBeTrue)
						})
					})
				})
			})

			Convey("And given the session was already committed", func() {
				customer := buildRegisteredCustomerWith(id)
				tx := test.BeginTx(db)
				session := repo.StartSession(tx)
				err := tx.Commit()
				So(err, ShouldBeNil)

				Convey("When the Customer is registered", func() {
					err = session.Register(customer)

					Convey("It should fail", func() {
						So(xerrors.Is(err, shared.ErrTechnical), ShouldBeTrue)
					})
				})
			})

			cleanUpArtefactsForEventSourcedRepositorySession(id)
		})

		Convey("And given an existing Customer", func() {
			id := values.GenerateID()
			customer := buildRegisteredCustomerWith(id)
			tx := test.BeginTx(db)
			session := repo.StartSession(tx)
			err := session.Register(customer)
			So(err, ShouldBeNil)
			err = tx.Commit()
			So(err, ShouldBeNil)

			Convey("When the same Customer is registered again", func() {
				customer := buildRegisteredCustomerWith(id)
				tx := test.BeginTx(db)
				session := repo.StartSession(tx)

				err = session.Register(customer)

				Convey("It should fail", func() {
					So(xerrors.Is(err, shared.ErrDuplicate), ShouldBeTrue)
				})
			})

			cleanUpArtefactsForEventSourcedRepositorySession(id)
		})
	})
}

func TestEventSourcedRepositorySession_Of(t *testing.T) {
	Convey("Given an existing Customer", t, func() {
		id := values.GenerateID()
		customer := buildRegisteredCustomerWith(id)
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		repo := diContainer.GetCustomerRepository()
		tx := test.BeginTx(db)
		session := repo.StartSession(tx)
		err := session.Register(customer)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)

		Convey("And given the Customer is not cached", func() {
			diContainer := test.SetUpDIContainer() // recreate the repository to reset the cache
			db := diContainer.GetPostgresDBConn()
			store := diContainer.GetPostgresEventStore()
			repo := diContainer.GetCustomerRepository()

			tx := test.BeginTx(db)

			Convey("When the Customer is retrieved", func() {
				session := repo.StartSession(tx)

				customer, err := session.Of(id)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)
					So(customer.ID().Equals(id), ShouldBeTrue)
					So(customer.StreamVersion(), ShouldEqual, uint(1))
				})

				err = tx.Commit()
				So(err, ShouldBeNil)

				Convey("And when the customer can not be reconstituted", func() {
					expectedErr := errors.New("mocked error")
					customerFactory := func(eventStream shared.DomainEvents) (domain.Customer, error) {
						return nil, expectedErr
					}

					repo := customers.NewEventSourcedRepository(store, customerFactory, customers.NewIdentityMap())
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
		})

		Convey("And given the Customer is cached", func() {
			tx := test.BeginTx(db)
			session := repo.StartSession(tx)

			_, err = session.Of(id)
			So(err, ShouldBeNil)

			Convey("And given the Customer was concurrently modified", func() {
				//otherRepo, _, _ := test.SetUpDIContainer()

				diContainer := test.SetUpDIContainer()
				otherRepo := diContainer.GetCustomerRepository()

				tx := test.BeginTx(db)

				session := otherRepo.StartSession(tx)

				sameCustomer, err := session.Of(id)
				So(err, ShouldBeNil)

				newEmailAddress := fmt.Sprintf("john+changed+%s@doe.com", id.String())
				changeEmailAddress, err := commands.NewChangeEmailAddress(id.String(), newEmailAddress)

				err = sameCustomer.ChangeEmailAddress(changeEmailAddress)
				So(err, ShouldBeNil)

				err = session.Persist(sameCustomer)
				So(err, ShouldBeNil)

				err = tx.Commit()
				So(err, ShouldBeNil)

				Convey("When the Customer is retrieved", func() {
					tx := test.BeginTx(db)
					session := repo.StartSession(tx)

					customer, err := session.Of(id)

					Convey("It should retrieve the up-to-date Customer", func() {
						So(err, ShouldBeNil)
						So(customer.ID().Equals(id), ShouldBeTrue)
						So(customer.StreamVersion(), ShouldEqual, uint(2))
					})

					err = tx.Commit()
					So(err, ShouldBeNil)

					Convey("And when the DB connection was closed", func() {
						tx := test.BeginTx(db)
						session := repo.StartSession(tx)

						err = db.Close()
						So(err, ShouldBeNil)

						customer, err := session.Of(id)

						Convey("It should fail", func() {
							So(xerrors.Is(err, shared.ErrTechnical), ShouldBeTrue)
							So(customer, ShouldBeNil)
						})
					})
				})
			})
		})

		cleanUpArtefactsForEventSourcedRepositorySession(id)
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

func TestEventSourcedRepositorySession_Persist(t *testing.T) {
	Convey("Given a changed Customer", t, func() {
		id := values.GenerateID()
		customer := buildRegisteredCustomerWith(id)
		diContainer := test.SetUpDIContainer()
		db := diContainer.GetPostgresDBConn()
		repo := diContainer.GetCustomerRepository()
		tx := test.BeginTx(db)
		session := repo.StartSession(tx)
		err := session.Register(customer)
		So(err, ShouldBeNil)
		err = tx.Commit()
		So(err, ShouldBeNil)
		changeEmailAddress, err := commands.NewChangeEmailAddress(
			id.String(),
			fmt.Sprintf("john+%s+changed@doe.com", id.String()),
		)
		So(err, ShouldBeNil)

		err = customer.ChangeEmailAddress(changeEmailAddress)
		So(err, ShouldBeNil)

		Convey("When the Customer is persisted", func() {
			tx := test.BeginTx(db)
			session := repo.StartSession(tx)

			err = session.Persist(customer)

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
				err = session.Persist(customer)

				Convey("It should fail", func() {
					So(xerrors.Is(err, shared.ErrTechnical), ShouldBeTrue)
				})
			})
		})

		cleanUpArtefactsForEventSourcedRepositorySession(id)
	})
}

/*** Test Helper Methods ***/

func buildRegisteredCustomerWith(id *values.ID) domain.Customer {
	emailAddress := fmt.Sprintf("john+%s@doe.com", id.String())
	givenName := "John"
	familyName := "Doe"
	register, err := commands.NewRegister(id.String(), emailAddress, givenName, familyName)
	So(err, ShouldBeNil)
	customer := domain.NewCustomerWith(register)

	return customer
}

//func startTxForPostgresEventSourcedRepositorySession(db *sql.DB) *sql.Tx {
//	tx, err := db.Begin()
//	So(err, ShouldBeNil)
//
//	return tx
//}

func cleanUpArtefactsForEventSourcedRepositorySession(id *values.ID) {
	diContainer := test.SetUpDIContainer()
	store := diContainer.GetPostgresEventStore()

	streamID := shared.NewStreamID("customer" + "-" + id.String())
	err := store.PurgeEventStream(streamID)
	So(err, ShouldBeNil)
}
