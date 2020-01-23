package application_test

import (
	"errors"
	"go-iddd/service/customer/application"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/customer/application/mocks"
	"go-iddd/service/lib"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func Test_CommandHandler_ConfirmEmailAddress_WithErrorFromCustomers(t *testing.T) {
	Convey("Setup", t, func() {
		customers := new(mocks.ForStoringCustomers)

		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		commandHandler := application.NewCommandHandler(customers, db)

		customerID := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("john@doe.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("John", "Doe")

		eventStream := lib.DomainEvents{
			events.ItWasRegistered(customerID, emailAddress, confirmationHash, personName, uint(1)),
		}

		confirmEmailAddress, err := commands.NewConfirmEmailAddress(
			customerID.ID(),
			emailAddress.EmailAddress(),
			confirmationHash.Hash(),
		)
		So(err, ShouldBeNil)

		Convey("Given the Customer can't be persisted because of a technical error", func() {
			customers.On("EventStream", mock.Anything).Return(eventStream, nil).Once()
			customers.On("Persist", mock.Anything, mock.Anything, mock.Anything).Return(lib.ErrTechnical).Once()

			Convey("When a Customer's emailAddress is confirmed", func() {
				err := commandHandler.ConfirmEmailAddress(confirmEmailAddress)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
					So(dbMock.ExpectationsWereMet(), ShouldBeNil)
				})
			})
		})
	})
}

func Test_CommandHandler_ChangeEmailAddress_WithErrorFromCustomers(t *testing.T) {
	Convey("Setup", t, func() {
		customers := new(mocks.ForStoringCustomers)

		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)
		dbMock.ExpectBegin()
		dbMock.ExpectRollback()

		commandHandler := application.NewCommandHandler(customers, db)

		customerID := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("john@doe.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("John", "Doe")

		eventStream := lib.DomainEvents{
			events.ItWasRegistered(customerID, emailAddress, confirmationHash, personName, uint(1)),
		}

		changeEmailAddress, err := commands.NewChangeEmailAddress(
			customerID.ID(),
			"john+changed@doe.com",
		)
		So(err, ShouldBeNil)

		Convey("Given the Customer can't be persisted because of a technical error", func() {
			customers.On("EventStream", mock.Anything).Return(eventStream, nil).Once()
			customers.On("Persist", mock.Anything, mock.Anything, mock.Anything).Return(lib.ErrTechnical).Once()

			Convey("When a Customer's emailAddress is confirmed", func() {
				err := commandHandler.ChangeEmailAddress(changeEmailAddress)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
					So(dbMock.ExpectationsWereMet(), ShouldBeNil)
				})
			})
		})
	})
}

func Test_CommandHandler_WithTransactionErrors(t *testing.T) {
	Convey("Setup", t, func() {
		customers := new(mocks.ForStoringCustomers)

		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)

		register, err := commands.NewRegister("john@doe.com", "John", "Doe")
		So(err, ShouldBeNil)

		commandHandler := application.NewCommandHandler(customers, db)

		Convey("Given begin transaction fails", func() {
			dbMock.ExpectBegin().WillReturnError(lib.ErrTechnical)

			Convey("When a command is handled", func() {
				err := commandHandler.Register(register)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
					So(dbMock.ExpectationsWereMet(), ShouldBeNil)
				})
			})
		})

		Convey("Given commit transaction fails", func() {
			dbMock.ExpectBegin()
			dbMock.ExpectCommit().WillReturnError(lib.ErrTechnical)

			customers.On("Register", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

			Convey("When a command is handled", func() {
				err := commandHandler.Register(register)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
					So(dbMock.ExpectationsWereMet(), ShouldBeNil)
				})
			})
		})
	})
}

func Test_CommandHandler_RetriesWithConcurrencyConflicts(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		customers := new(mocks.ForStoringCustomers)

		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)

		commandHandler := application.NewCommandHandler(customers, db)

		customerID := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("john@doe.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("John", "Doe")
		eventStream := lib.DomainEvents{
			events.ItWasRegistered(customerID, emailAddress, confirmationHash, personName, uint(1)),
		}

		Convey("And given a ChangeEmailAddress command", func() {
			changeEmailAddress, err := commands.NewChangeEmailAddress(
				customerID.ID(),
				emailAddress.EmailAddress(),
			)
			So(err, ShouldBeNil)

			Convey("When the command is handled", func() {
				Convey("And when finding the Customer succeeds", func() {

					Convey("And when saving the Customer has a concurrency conflict once", func() {
						// should be called twice due to retry
						customers.On("EventStream", mock.Anything).Return(eventStream, nil).Twice()

						// fist attempt runs into a concurrency conflict
						customers.
							On("Persist", mock.Anything, mock.Anything, mock.Anything).
							Return(lib.ErrConcurrencyConflict).
							Once()
						dbMock.ExpectBegin()
						dbMock.ExpectRollback()

						// second attempt works
						customers.
							On("Persist", mock.Anything, mock.Anything, mock.Anything).
							Return(nil).
							Once()
						dbMock.ExpectBegin()
						dbMock.ExpectCommit()

						err := commandHandler.ChangeEmailAddress(changeEmailAddress)

						Convey("It should modify the Customer and save it", func() {
							So(err, ShouldBeNil)
							So(dbMock.ExpectationsWereMet(), ShouldBeNil)
						})
					})

					Convey("And when saving the Customer has too many concurrency conflicts", func() {
						// should be called 10 times due to retries
						customers.
							On("EventStream", mock.Anything).
							Return(eventStream, nil).
							Times(10)

						// all attempts run into a concurrency conflict
						customers.
							On("Persist", mock.Anything, mock.Anything, mock.Anything).
							Return(lib.ErrConcurrencyConflict).
							Times(10)

						// we expect this 10 time - no simpler way with Sqlmock
						for i := 1; i <= 10; i++ {
							dbMock.ExpectBegin()
							dbMock.ExpectRollback()
						}

						err := commandHandler.ChangeEmailAddress(changeEmailAddress)

						Convey("It should fail", func() {
							So(err, ShouldNotBeNil)
							So(errors.Is(err, lib.ErrConcurrencyConflict), ShouldBeTrue)
							So(dbMock.ExpectationsWereMet(), ShouldBeNil)
						})
					})
				})
			})
		})
	})
}
