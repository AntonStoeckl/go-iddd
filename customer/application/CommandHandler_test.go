package application_test

import (
	"go-iddd/customer/application"
	"go-iddd/customer/application/mocks"
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

/*** Test factory method ***/

func TestNewCommandHandler(t *testing.T) {
	Convey("When a new CommandHandler is created", t, func() {
		sessionStarter := new(mocks.StartsCustomersSession)
		db, _, err := sqlmock.New()
		So(err, ShouldBeNil)

		commandHandler := application.NewCommandHandler(sessionStarter, db)

		Convey("It should succeed", func() {
			So(commandHandler, ShouldNotBeNil)
			So(commandHandler, ShouldHaveSameTypeAs, (*application.CommandHandler)(nil))
		})
	})
}

/*** Test business cases ***/

func TestCommandHandler_Handle_Register(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		customerID := values.GenerateCustomerID()
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
				customerID.ID(),
				emailAddress.EmailAddress(),
				"John",
				"Doe",
			)
			So(err, ShouldBeNil)

			Convey("When the command is handled", func() {
				Convey("And when registering the Customer succeeds", func() {
					customers.On("Register", customerID, mock.AnythingOfType("shared.DomainEvents")).Return(nil).Once()
					dbMock.ExpectCommit()

					err := commandHandler.Register(register)

					Convey("It should register and save a Customer", func() {
						So(err, ShouldBeNil)
						So(customers.AssertExpectations(t), ShouldBeTrue)
					})
				})

				Convey("And when registering the Customer fails", func() {
					customers.On("Register", customerID, mock.AnythingOfType("shared.DomainEvents")).Return(shared.ErrTechnical).Once()
					dbMock.ExpectRollback()

					err := commandHandler.Register(register)

					Convey("It should fail", func() {
						So(errors.Is(err, shared.ErrTechnical), ShouldBeTrue)
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
		emailAddress, err := values.BuildEmailAddress("john@doe.com")
		So(err, ShouldBeNil)

		recordedEvents := registerCustomerForCommandHandlerTest(customerID, emailAddress)
		confirmationHash := recordedEvents[0].(events.Registered).ConfirmationHash()

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
					customers.On("EventStream", confirmEmailAddress.AggregateID()).Return(recordedEvents, nil).Once()

					Convey("And when executing the command succeeds", func() {
						Convey("And when persisting the Customer succeeds", func() {
							customers.On("Persist", customerID, mock.AnythingOfType("shared.DomainEvents")).Return(nil).Once()
							dbMock.ExpectCommit()

							err := commandHandler.ConfirmEmailAddress(confirmEmailAddress)

							Convey("It should modify the Customer and save it", func() {
								So(err, ShouldBeNil)
								So(customers.AssertExpectations(t), ShouldBeTrue)
								So(dbMock.ExpectationsWereMet(), ShouldBeNil)
							})
						})

						Convey("And when persisting the Customer fails", func() {
							customers.On("Persist", customerID, mock.AnythingOfType("shared.DomainEvents")).Return(shared.ErrTechnical).Once()
							dbMock.ExpectRollback()

							err := commandHandler.ConfirmEmailAddress(confirmEmailAddress)

							Convey("It should fail", func() {
								So(errors.Is(err, shared.ErrTechnical), ShouldBeTrue)
								So(customers.AssertExpectations(t), ShouldBeTrue)
								So(dbMock.ExpectationsWereMet(), ShouldBeNil)
							})
						})
					})

					Convey("And when executing the command fails", func() {
						customers.On("Persist", customerID, mock.AnythingOfType("shared.DomainEvents")).Return(nil).Once()
						dbMock.ExpectCommit()

						confirmEmailAddress, err := commands.NewConfirmEmailAddress(
							customerID.ID(),
							emailAddress.EmailAddress(),
							"invalid_hash",
						)
						So(err, ShouldBeNil)

						err = commandHandler.ConfirmEmailAddress(confirmEmailAddress)

						Convey("It should fail", func() {
							So(errors.Is(err, shared.ErrDomainConstraintsViolation), ShouldBeTrue)
							So(customers.AssertExpectations(t), ShouldBeTrue)
							So(dbMock.ExpectationsWereMet(), ShouldBeNil)
						})
					})
				})

				Convey("And when finding the Customer fails", func() {
					customers.On("EventStream", confirmEmailAddress.AggregateID()).Return(nil, shared.ErrTechnical).Once()
					dbMock.ExpectRollback()

					err := commandHandler.ConfirmEmailAddress(confirmEmailAddress)

					Convey("It should fail", func() {
						So(errors.Is(err, shared.ErrTechnical), ShouldBeTrue)
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
		emailAddress, err := values.BuildEmailAddress("john@doe.com")
		So(err, ShouldBeNil)

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
					recordedEvents := registerCustomerForCommandHandlerTest(customerID, emailAddress)
					customers.On("EventStream", changeEmailAddress.AggregateID()).Return(recordedEvents, nil).Once()

					Convey("And when saving the Customer succeeds", func() {
						customers.On("Persist", customerID, mock.AnythingOfType("shared.DomainEvents")).Return(nil).Once()
						dbMock.ExpectCommit()

						err := commandHandler.ChangeEmailAddress(changeEmailAddress)

						Convey("It should modify the Customer and save it", func() {
							So(err, ShouldBeNil)
							So(customers.AssertExpectations(t), ShouldBeTrue)
							So(dbMock.ExpectationsWereMet(), ShouldBeNil)
						})
					})

					Convey("And when saving the Customer fails", func() {
						customers.On("Persist", customerID, mock.AnythingOfType("shared.DomainEvents")).Return(shared.ErrTechnical).Once()
						dbMock.ExpectRollback()

						err := commandHandler.ChangeEmailAddress(changeEmailAddress)

						Convey("It should fail", func() {
							So(errors.Is(err, shared.ErrTechnical), ShouldBeTrue)
							So(customers.AssertExpectations(t), ShouldBeTrue)
							So(dbMock.ExpectationsWereMet(), ShouldBeNil)
						})
					})
				})

				Convey("And when finding the Customer fails", func() {
					customers.On("EventStream", changeEmailAddress.AggregateID()).Return(nil, shared.ErrTechnical).Once()
					dbMock.ExpectRollback()

					err := commandHandler.ChangeEmailAddress(changeEmailAddress)

					Convey("It should fail", func() {
						So(errors.Is(err, shared.ErrTechnical), ShouldBeTrue)
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
		emailAddress, err := values.BuildEmailAddress("john@doe.com")
		So(err, ShouldBeNil)

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
					recordedEvents := registerCustomerForCommandHandlerTest(customerID, emailAddress)

					Convey("And when saving the Customer has a concurrency conflict once", func() {
						// should be called twice due to retry
						customers.On("EventStream", changeEmailAddress.AggregateID()).Return(recordedEvents, nil).Twice()

						// fist attempt runs into a concurrency conflict
						customers.On("Persist", customerID, mock.AnythingOfType("shared.DomainEvents")).Return(shared.ErrConcurrencyConflict).Once()
						dbMock.ExpectBegin()
						dbMock.ExpectRollback()

						// second attempt works
						customers.On("Persist", customerID, mock.AnythingOfType("shared.DomainEvents")).Return(nil).Once()
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
						customers.On("EventStream", changeEmailAddress.AggregateID()).Return(recordedEvents, nil).Times(10)

						// all attempts run into a concurrency conflict
						customers.On("Persist", customerID, mock.AnythingOfType("shared.DomainEvents")).Return(shared.ErrConcurrencyConflict).Times(10)

						// we expect this 10 time - no simpler way with Sqlmock
						for i := 1; i <= 10; i++ {
							dbMock.ExpectBegin()
							dbMock.ExpectRollback()
						}

						err := commandHandler.ChangeEmailAddress(changeEmailAddress)

						Convey("It should fail", func() {
							So(err, ShouldNotBeNil)
							So(errors.Is(err, shared.ErrConcurrencyConflict), ShouldBeTrue)
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
				So(errors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
			})
		})

		// Convey("When an unknown command is handled", func() {
		// 	unknownCommand := new(mocks.Command)
		// 	unknownCommand.On("AggregateID").Return(values.GenerateCustomerID())
		// 	unknownCommand.On("CommandName").Return("unknown")
		// 	unknownCommand.On("ShouldBeValid").Return(nil)
		//
		// 	err := commandHandler.Handle(unknownCommand)
		//
		// 	Convey("It should fail", func() {
		// 		So(err, ShouldBeError)
		// 		So(errors.Is(err, shared.ErrCommandIsUnknown), ShouldBeTrue)
		// 	})
		// })
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
			sqlMock.ExpectBegin().WillReturnError(shared.ErrTechnical)

			register, err := commands.NewRegister(
				"64bcf656-da30-4f5a-b0b5-aead60965aa3",
				"john@doe.com",
				"John",
				"Doe",
			)
			So(err, ShouldBeNil)

			err = commandHandler.Register(register)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, shared.ErrTechnical), ShouldBeTrue)
				So(sessionStarter.AssertExpectations(t), ShouldBeTrue)
				So(customers.AssertExpectations(t), ShouldBeTrue)
			})
		})

		Convey("When committing the repositry session fails", func() {
			sqlMock.ExpectBegin()
			sessionStarter.On("StartSession", mock.AnythingOfType("*sql.Tx")).Return(customers)
			customers.
				On("Register", mock.AnythingOfType("values.CustomerID"), mock.AnythingOfType("shared.DomainEvents")).
				Return(nil)
			sqlMock.ExpectCommit().WillReturnError(shared.ErrTechnical)

			register, err := commands.NewRegister(
				"64bcf656-da30-4f5a-b0b5-aead60965aa3",
				"john@doe.com",
				"John",
				"Doe",
			)
			So(err, ShouldBeNil)

			err = commandHandler.Register(register)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, shared.ErrTechnical), ShouldBeTrue)
				// So(sessionStarter.AssertExpectations(t), ShouldBeTrue)
				So(customers.AssertExpectations(t), ShouldBeTrue)
			})
		})
	})
}

func registerCustomerForCommandHandlerTest(id values.CustomerID, emailAddress values.EmailAddress) shared.DomainEvents {
	register, err := commands.NewRegister(
		id.ID(),
		emailAddress.EmailAddress(),
		"John",
		"Doe",
	)
	So(err, ShouldBeNil)

	recordedEvents := domain.RegisterCustomer(register)
	So(recordedEvents, ShouldHaveLength, 1)
	So(recordedEvents[0], ShouldHaveSameTypeAs, events.Registered{})

	return recordedEvents
}
