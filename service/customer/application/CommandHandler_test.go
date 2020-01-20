package application_test

import (
	"go-iddd/service/customer/application"
	"go-iddd/service/customer/application/mocks"
	"go-iddd/service/customer/domain/commands"
	"go-iddd/service/customer/domain/events"
	"go-iddd/service/customer/domain/values"
	"go-iddd/service/lib"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

/*** Test business cases ***/

func TestCommandHandler_Handle_Register(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		emailAddress, err := values.BuildEmailAddress("john@doe.com")
		So(err, ShouldBeNil)

		customers := new(mocks.Customers)

		sessionStarter := new(mocks.StartsCustomersSession)
		sessionStarter.On("StartSession", mock.AnythingOfType("*sql.Tx")).Return(customers)

		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)
		dbMock.ExpectBegin()

		commandHandler := application.NewCommandHandler(sessionStarter, db)

		Convey("And given a Register command", func() {
			register, err := commands.NewRegister(
				emailAddress.EmailAddress(),
				"John",
				"Doe",
			)
			So(err, ShouldBeNil)

			customerID := register.CustomerID()

			Convey("When the command is handled", func() {
				Convey("And when registering the Customer succeeds", func() {
					customers.On("Register", customerID, mock.AnythingOfType("lib.DomainEvents")).Return(nil).Once()
					dbMock.ExpectCommit()

					err := commandHandler.Register(register)

					Convey("It should succeed", func() {
						So(err, ShouldBeNil)
						So(customers.AssertExpectations(t), ShouldBeTrue)
					})
				})

				Convey("And when registering the Customer fails", func() {
					customers.On("Register", customerID, mock.AnythingOfType("lib.DomainEvents")).Return(lib.ErrTechnical).Once()
					dbMock.ExpectRollback()

					err := commandHandler.Register(register)

					Convey("It should fail", func() {
						So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						So(customers.AssertExpectations(t), ShouldBeTrue)
					})
				})
			})
		})
	})
}

func TestCommandHandler_Handle_ConfirmEmailAddress(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		customerID := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("john@doe.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("John", "Doe")
		registered := events.ItWasRegistered(customerID, emailAddress, confirmationHash, personName, uint(1))
		eventStream := lib.DomainEvents{registered}

		customers := new(mocks.Customers)

		sessionStarter := new(mocks.StartsCustomersSession)
		sessionStarter.On("StartSession", mock.AnythingOfType("*sql.Tx")).Return(customers)

		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)
		dbMock.ExpectBegin()

		commandHandler := application.NewCommandHandler(sessionStarter, db)

		Convey("And given a ConfirmEmailAddress command", func() {
			confirmEmailAddress, err := commands.NewConfirmEmailAddress(
				customerID.ID(),
				emailAddress.EmailAddress(),
				confirmationHash.Hash(),
			)
			So(err, ShouldBeNil)

			Convey("When the command is handled", func() {
				Convey("And when finding the Customer succeeds", func() {
					customers.On("EventStream", confirmEmailAddress.AggregateID()).Return(eventStream, nil).Once()

					Convey("And when executing the command succeeds", func() {
						Convey("And when persisting the Customer succeeds", func() {
							customers.On("Persist", customerID, mock.AnythingOfType("lib.DomainEvents")).Return(nil).Once()
							dbMock.ExpectCommit()

							err := commandHandler.ConfirmEmailAddress(confirmEmailAddress)

							Convey("It should modify the Customer and save it", func() {
								So(err, ShouldBeNil)
								So(customers.AssertExpectations(t), ShouldBeTrue)
								So(dbMock.ExpectationsWereMet(), ShouldBeNil)
							})
						})

						Convey("And when persisting the Customer fails", func() {
							customers.On("Persist", customerID, mock.AnythingOfType("lib.DomainEvents")).Return(lib.ErrTechnical).Once()
							dbMock.ExpectRollback()

							err := commandHandler.ConfirmEmailAddress(confirmEmailAddress)

							Convey("It should fail", func() {
								So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
								So(customers.AssertExpectations(t), ShouldBeTrue)
								So(dbMock.ExpectationsWereMet(), ShouldBeNil)
							})
						})
					})

					Convey("And when executing the command fails", func() {
						customers.On("Persist", customerID, mock.AnythingOfType("lib.DomainEvents")).Return(nil).Once()
						dbMock.ExpectCommit()

						confirmEmailAddress, err := commands.NewConfirmEmailAddress(
							customerID.ID(),
							emailAddress.EmailAddress(),
							"invalid_hash",
						)
						So(err, ShouldBeNil)

						err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)

						Convey("It should fail", func() {
							So(errors.Is(err, lib.ErrDomainConstraintsViolation), ShouldBeTrue)
							So(customers.AssertExpectations(t), ShouldBeTrue)
							So(dbMock.ExpectationsWereMet(), ShouldBeNil)
						})
					})
				})

				Convey("And when finding the Customer fails", func() {
					customers.On("EventStream", confirmEmailAddress.AggregateID()).Return(nil, lib.ErrTechnical).Once()
					dbMock.ExpectRollback()

					err := commandHandler.ConfirmEmailAddress(confirmEmailAddress)

					Convey("It should fail", func() {
						So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						So(customers.AssertExpectations(t), ShouldBeTrue)
						So(dbMock.ExpectationsWereMet(), ShouldBeNil)
					})
				})
			})
		})
	})
}

func TestCommandHandler_Handle_ChangeEmailAddress(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		customerID := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("john@doe.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("John", "Doe")
		registered := events.ItWasRegistered(customerID, emailAddress, confirmationHash, personName, uint(1))
		eventStream := lib.DomainEvents{registered}

		customers := new(mocks.Customers)

		sessionStarter := new(mocks.StartsCustomersSession)
		sessionStarter.On("StartSession", mock.AnythingOfType("*sql.Tx")).Return(customers)

		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)
		dbMock.ExpectBegin()

		commandHandler := application.NewCommandHandler(sessionStarter, db)

		Convey("And given a ChangeEmailAddress command", func() {
			changeEmailAddress, err := commands.NewChangeEmailAddress(
				customerID.ID(),
				emailAddress.EmailAddress(),
			)
			So(err, ShouldBeNil)

			Convey("When the command is handled", func() {
				Convey("And when finding the Customer succeeds", func() {
					customers.On("EventStream", changeEmailAddress.AggregateID()).Return(eventStream, nil).Once()

					Convey("And when saving the Customer succeeds", func() {
						customers.On("Persist", customerID, mock.AnythingOfType("lib.DomainEvents")).Return(nil).Once()
						dbMock.ExpectCommit()

						err := commandHandler.ChangeEmailAddress(changeEmailAddress)

						Convey("It should modify the Customer and save it", func() {
							So(err, ShouldBeNil)
							So(customers.AssertExpectations(t), ShouldBeTrue)
							So(dbMock.ExpectationsWereMet(), ShouldBeNil)
						})
					})

					Convey("And when saving the Customer fails", func() {
						customers.On("Persist", customerID, mock.AnythingOfType("lib.DomainEvents")).Return(lib.ErrTechnical).Once()
						dbMock.ExpectRollback()

						err := commandHandler.ChangeEmailAddress(changeEmailAddress)

						Convey("It should fail", func() {
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
							So(customers.AssertExpectations(t), ShouldBeTrue)
							So(dbMock.ExpectationsWereMet(), ShouldBeNil)
						})
					})
				})

				Convey("And when finding the Customer fails", func() {
					customers.On("EventStream", changeEmailAddress.AggregateID()).Return(nil, lib.ErrTechnical).Once()
					dbMock.ExpectRollback()

					err := commandHandler.ChangeEmailAddress(changeEmailAddress)

					Convey("It should fail", func() {
						So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						So(customers.AssertExpectations(t), ShouldBeTrue)
						So(dbMock.ExpectationsWereMet(), ShouldBeNil)
					})
				})
			})
		})
	})
}

/*** Test generic error cases ***/

func TestCommandHandler_Handle_RetriesWithConcurrencyConflicts(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		customerID := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("john@doe.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("John", "Doe")
		registered := events.ItWasRegistered(customerID, emailAddress, confirmationHash, personName, uint(1))
		eventStream := lib.DomainEvents{registered}

		customers := new(mocks.Customers)

		sessionStarter := new(mocks.StartsCustomersSession)
		sessionStarter.On("StartSession", mock.AnythingOfType("*sql.Tx")).Return(customers)

		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)

		commandHandler := application.NewCommandHandler(sessionStarter, db)

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
						customers.On("EventStream", changeEmailAddress.AggregateID()).Return(eventStream, nil).Twice()

						// fist attempt runs into a concurrency conflict
						customers.On("Persist", customerID, mock.AnythingOfType("lib.DomainEvents")).Return(lib.ErrConcurrencyConflict).Once()
						dbMock.ExpectBegin()
						dbMock.ExpectRollback()

						// second attempt works
						customers.On("Persist", customerID, mock.AnythingOfType("lib.DomainEvents")).Return(nil).Once()
						dbMock.ExpectBegin()
						dbMock.ExpectCommit()

						err := commandHandler.ChangeEmailAddress(changeEmailAddress)

						Convey("It should modify the Customer and save it", func() {
							So(err, ShouldBeNil)
							So(customers.AssertExpectations(t), ShouldBeTrue)
							So(dbMock.ExpectationsWereMet(), ShouldBeNil)
						})
					})

					Convey("And when saving the Customer has too many concurrency conflicts", func() {
						// should be called 10 times due to retries
						customers.On("EventStream", changeEmailAddress.AggregateID()).Return(eventStream, nil).Times(10)

						// all attempts run into a concurrency conflict
						customers.On("Persist", customerID, mock.AnythingOfType("lib.DomainEvents")).Return(lib.ErrConcurrencyConflict).Times(10)

						// we expect this 10 time - no simpler way with Sqlmock
						for i := 1; i <= 10; i++ {
							dbMock.ExpectBegin()
							dbMock.ExpectRollback()
						}

						err := commandHandler.ChangeEmailAddress(changeEmailAddress)

						Convey("It should fail", func() {
							So(err, ShouldNotBeNil)
							So(errors.Is(err, lib.ErrConcurrencyConflict), ShouldBeTrue)
							So(customers.AssertExpectations(t), ShouldBeTrue)
							So(dbMock.ExpectationsWereMet(), ShouldBeNil)
						})
					})
				})
			})
		})
	})
}

func TestCommandHandler_Handle_WithInvalidCommand(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		customers := new(mocks.Customers)

		sessionStarter := new(mocks.StartsCustomersSession)
		sessionStarter.On("StartSession", mock.AnythingOfType("*sql.Tx")).Return(customers)

		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)
		dbMock.ExpectBegin()

		commandHandler := application.NewCommandHandler(sessionStarter, db)

		Convey("When an empty command is handled", func() {
			emptyCommand := commands.ConfirmEmailAddress{}

			err := commandHandler.ConfirmEmailAddress(emptyCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
			})
		})
	})
}

func TestCommandHandler_Handle_WithSessionErrors(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		db, sqlMock, err := sqlmock.New()
		So(err, ShouldBeNil)
		sessionStarter := new(mocks.StartsCustomersSession)
		customers := new(mocks.Customers)

		commandHandler := application.NewCommandHandler(sessionStarter, db)

		Convey("When beginning the transaction fails", func() {
			sqlMock.ExpectBegin().WillReturnError(lib.ErrTechnical)

			register, err := commands.NewRegister(
				"john@doe.com",
				"John",
				"Doe",
			)
			So(err, ShouldBeNil)

			err = commandHandler.Register(register)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				So(sessionStarter.AssertExpectations(t), ShouldBeTrue)
				So(customers.AssertExpectations(t), ShouldBeTrue)
			})
		})

		Convey("When committing the repositry session fails", func() {
			sqlMock.ExpectBegin()
			sessionStarter.On("StartSession", mock.AnythingOfType("*sql.Tx")).Return(customers)
			customers.
				On("Register", mock.AnythingOfType("values.CustomerID"), mock.AnythingOfType("lib.DomainEvents")).
				Return(nil)
			sqlMock.ExpectCommit().WillReturnError(lib.ErrTechnical)

			register, err := commands.NewRegister(
				"john@doe.com",
				"John",
				"Doe",
			)
			So(err, ShouldBeNil)

			err = commandHandler.Register(register)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				// So(sessionStarter.AssertExpectations(t), ShouldBeTrue)
				So(customers.AssertExpectations(t), ShouldBeTrue)
			})
		})
	})
}
