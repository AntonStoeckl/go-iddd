package application_test

import (
	"fmt"
	"go-iddd/customer/application"
	"go-iddd/customer/application/mocks"
	"go-iddd/customer/domain/commands"
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
		repo := new(mocks.StartsRepositorySessions)
		db, _, err := sqlmock.New()
		So(err, ShouldBeNil)

		commandHandler := application.NewCommandHandler(repo, db)

		Convey("It should succeed", func() {
			So(commandHandler, ShouldNotBeNil)
			So(commandHandler, ShouldImplement, (*shared.CommandHandler)(nil))
			So(commandHandler, ShouldHaveSameTypeAs, (*application.CommandHandler)(nil))
		})
	})
}

/*** Test business cases ***/

func TestCommandHandler_Handle_Register(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		customers := new(mocks.PersistableCustomers)

		repo := new(mocks.StartsRepositorySessions)
		repo.On("StartSession", mock.AnythingOfType("*sql.Tx")).Return(customers)

		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)
		dbMock.ExpectBegin()

		commandHandler := application.NewCommandHandler(repo, db)

		Convey("And given a Register command", func() {
			register, err := commands.NewRegister(
				"64bcf656-da30-4f5a-b0b5-aead60965aa3",
				"john@doe.com",
				"John",
				"Doe",
			)
			So(err, ShouldBeNil)

			Convey("When the command is handled", func() {
				Convey("And when saving the Customer succeeds", func() {
					customers.On("Register", mock.AnythingOfType("*domain.customer")).Return(nil).Once()
					dbMock.ExpectCommit()

					err := commandHandler.Handle(register)

					Convey("It should register and save a Customer", func() {
						So(err, ShouldBeNil)
						So(customers.AssertExpectations(t), ShouldBeTrue)
					})
				})

				Convey("And when saving the Customer fails", func() {
					customers.On("Register", mock.AnythingOfType("*domain.customer")).Return(shared.ErrTechnical).Once()
					dbMock.ExpectRollback()

					err := commandHandler.Handle(register)

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
		customers := new(mocks.PersistableCustomers)

		repo := new(mocks.StartsRepositorySessions)
		repo.On("StartSession", mock.AnythingOfType("*sql.Tx")).Return(customers)

		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)
		dbMock.ExpectBegin()

		commandHandler := application.NewCommandHandler(repo, db)

		Convey("And given a ConfirmEmailAddress command", func() {
			confirmEmailAddress, err := commands.NewConfirmEmailAddress(
				"64bcf656-da30-4f5a-b0b5-aead60965aa3",
				"john@doe.com",
				"secret_hash",
			)
			So(err, ShouldBeNil)

			Convey("When the command is handled", func() {
				Convey("And when finding the Customer succeeds", func() {
					mockCustomer := new(mocks.Customer)
					customers.On("Of", confirmEmailAddress.AggregateID()).Return(mockCustomer, nil).Once()

					Convey("And when executing the command succeeds", func() {
						mockCustomer.On("ConfirmEmailAddress", confirmEmailAddress).Return(nil)

						Convey("And when saving the Customer succeeds", func() {
							customers.On("Persist", mockCustomer).Return(nil).Once()
							dbMock.ExpectCommit()

							err := commandHandler.Handle(confirmEmailAddress)

							Convey("It should modify the Customer and save it", func() {
								So(err, ShouldBeNil)
								So(mockCustomer.AssertExpectations(t), ShouldBeTrue)
								So(customers.AssertExpectations(t), ShouldBeTrue)
							})
						})

						Convey("And when saving the Customer fails", func() {
							customers.On("Persist", mockCustomer).Return(shared.ErrTechnical).Once()
							dbMock.ExpectRollback()

							err := commandHandler.Handle(confirmEmailAddress)

							Convey("It should fail", func() {
								So(errors.Is(err, shared.ErrTechnical), ShouldBeTrue)
								So(customers.AssertExpectations(t), ShouldBeTrue)
							})
						})
					})

					Convey("And when executing the command fails", func() {
						mockCustomer.On("ConfirmEmailAddress", confirmEmailAddress).Return(shared.ErrDomainConstraintsViolation)
						dbMock.ExpectRollback()

						err := commandHandler.Handle(confirmEmailAddress)

						Convey("It should fail", func() {
							So(errors.Is(err, shared.ErrDomainConstraintsViolation), ShouldBeTrue)
							So(mockCustomer.AssertExpectations(t), ShouldBeTrue)
							So(customers.AssertExpectations(t), ShouldBeTrue)
						})
					})
				})

				Convey("And when finding the Customer fails", func() {
					customers.On("Of", confirmEmailAddress.AggregateID()).Return(nil, shared.ErrTechnical).Once()
					dbMock.ExpectRollback()

					err := commandHandler.Handle(confirmEmailAddress)

					Convey("It should fail", func() {
						So(errors.Is(err, shared.ErrTechnical), ShouldBeTrue)
						So(customers.AssertExpectations(t), ShouldBeTrue)
					})
				})
			})
		})
	})
}

func TestCommandHandler_Handle_ChangeEmailAddress(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		customers := new(mocks.PersistableCustomers)

		repo := new(mocks.StartsRepositorySessions)
		repo.On("StartSession", mock.AnythingOfType("*sql.Tx")).Return(customers)

		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)
		dbMock.ExpectBegin()

		commandHandler := application.NewCommandHandler(repo, db)

		Convey("And given a ChangeEmailAddress command", func() {
			changeEmailAddress, err := commands.NewChangeEmailAddress(
				"64bcf656-da30-4f5a-b0b5-aead60965aa3",
				"john@doe.com",
			)
			So(err, ShouldBeNil)

			Convey("When the command is handled", func() {
				Convey("And when finding the Customer succeeds", func() {
					mockCustomer := new(mocks.Customer)
					customers.On("Of", changeEmailAddress.AggregateID()).Return(mockCustomer, nil).Once()

					Convey("And when executing the command succeeds", func() {
						mockCustomer.On("ChangeEmailAddress", changeEmailAddress).Return(nil)

						Convey("And when saving the Customer succeeds", func() {
							customers.On("Persist", mockCustomer).Return(nil).Once()
							dbMock.ExpectCommit()

							err := commandHandler.Handle(changeEmailAddress)

							Convey("It should modify the Customer and save it", func() {
								So(err, ShouldBeNil)
								So(mockCustomer.AssertExpectations(t), ShouldBeTrue)
								So(customers.AssertExpectations(t), ShouldBeTrue)
							})
						})

						Convey("And when saving the Customer fails", func() {
							customers.On("Persist", mockCustomer).Return(shared.ErrTechnical).Once()
							dbMock.ExpectRollback()

							err := commandHandler.Handle(changeEmailAddress)

							Convey("It should fail", func() {
								So(errors.Is(err, shared.ErrTechnical), ShouldBeTrue)
								So(customers.AssertExpectations(t), ShouldBeTrue)
							})
						})
					})

					Convey("And when executing the command fails", func() {
						mockCustomer.On("ChangeEmailAddress", changeEmailAddress).Return(shared.ErrDomainConstraintsViolation)
						dbMock.ExpectRollback()

						err := commandHandler.Handle(changeEmailAddress)

						Convey("It should fail", func() {
							So(errors.Is(err, shared.ErrDomainConstraintsViolation), ShouldBeTrue)
							So(mockCustomer.AssertExpectations(t), ShouldBeTrue)
							So(customers.AssertExpectations(t), ShouldBeTrue)
						})
					})
				})

				Convey("And when finding the Customer fails", func() {
					customers.On("Of", changeEmailAddress.AggregateID()).Return(nil, shared.ErrTechnical).Once()
					dbMock.ExpectRollback()

					err := commandHandler.Handle(changeEmailAddress)

					Convey("It should fail", func() {
						So(errors.Is(err, shared.ErrTechnical), ShouldBeTrue)
						So(customers.AssertExpectations(t), ShouldBeTrue)
					})
				})
			})
		})
	})
}

/*** Test generic error cases ***/

func TestCommandHandler_Handle_RetriesWithConcurrencyConflicts(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		customers := new(mocks.PersistableCustomers)

		repo := new(mocks.StartsRepositorySessions)
		repo.On("StartSession", mock.AnythingOfType("*sql.Tx")).Return(customers)

		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)

		commandHandler := application.NewCommandHandler(repo, db)

		Convey("And given a ChangeEmailAddress command", func() {
			changeEmailAddress, err := commands.NewChangeEmailAddress(
				"64bcf656-da30-4f5a-b0b5-aead60965aa3",
				"john@doe.com",
			)
			So(err, ShouldBeNil)

			Convey("When the command is handled", func() {
				Convey("And when finding the Customer succeeds", func() {
					mockCustomer := new(mocks.Customer)

					Convey("And when executing the command succeeds", func() {
						mockCustomer.On("ChangeEmailAddress", changeEmailAddress).Return(nil)

						Convey("And when saving the Customer has a concurrency conflict once", func() {
							// should be called twice due to retry
							customers.On("Of", changeEmailAddress.AggregateID()).Return(mockCustomer, nil).Twice()

							// fist attempt runs into a concurrency conflict
							customers.On("Persist", mockCustomer).Return(shared.ErrConcurrencyConflict).Once()
							dbMock.ExpectBegin()
							dbMock.ExpectRollback()

							// second attempt works
							customers.On("Persist", mockCustomer).Return(nil).Once()
							dbMock.ExpectBegin()
							dbMock.ExpectCommit()

							err := commandHandler.Handle(changeEmailAddress)

							Convey("It should modify the Customer and save it", func() {
								So(err, ShouldBeNil)
								So(mockCustomer.AssertExpectations(t), ShouldBeTrue)
								So(customers.AssertExpectations(t), ShouldBeTrue)
							})
						})

						Convey("And when saving the Customer has too many concurrency conflicts", func() {
							// should be called 10 times due to retries
							customers.On("Of", changeEmailAddress.AggregateID()).Return(mockCustomer, nil).Times(10)

							// all attempts run into a concurrency conflict
							customers.On("Persist", mockCustomer).Return(shared.ErrConcurrencyConflict).Times(10)

							// we expect this 10 time - no simpler way with Sqlmock
							for i := 1; i <= 10; i++ {
								dbMock.ExpectBegin()
								dbMock.ExpectRollback()
							}

							err := commandHandler.Handle(changeEmailAddress)

							Convey("It should fail", func() {
								So(err, ShouldNotBeNil)
								So(errors.Is(err, shared.ErrConcurrencyConflict), ShouldBeTrue)
								So(mockCustomer.AssertExpectations(t), ShouldBeTrue)
								So(customers.AssertExpectations(t), ShouldBeTrue)
							})
						})
					})
				})
			})
		})
	})
}

func TestCommandHandler_Handle_WithInvalidCommand(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		customers := new(mocks.PersistableCustomers)

		repo := new(mocks.StartsRepositorySessions)
		repo.On("StartSession", mock.AnythingOfType("*sql.Tx")).Return(customers)

		db, dbMock, err := sqlmock.New()
		So(err, ShouldBeNil)
		dbMock.ExpectBegin()

		commandHandler := application.NewCommandHandler(repo, db)

		Convey("When a nil interface command is handled", func() {
			var nilInterfaceCommand shared.Command
			err := commandHandler.Handle(nilInterfaceCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
				fmt.Println(err)
			})
		})

		Convey("When a nil pointer command is handled", func() {
			var nilCommand *commands.ConfirmEmailAddress
			err := commandHandler.Handle(nilCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
				fmt.Println(err)
			})
		})

		Convey("When an empty command is handled", func() {
			emptyCommand := &commands.ConfirmEmailAddress{}
			err := commandHandler.Handle(emptyCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
				fmt.Println(err)
			})
		})

		Convey("When an unknown command is handled", func() {
			unknownCommand := new(mocks.Command)
			unknownCommand.On("AggregateID").Return(values.GenerateID())
			unknownCommand.On("CommandName").Return("unknown")

			err := commandHandler.Handle(unknownCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(errors.Is(err, shared.ErrCommandIsUnknown), ShouldBeTrue)
				fmt.Println(err)
			})
		})
	})
}

//func TestCommandHandler_Handle_WithSessionErrors(t *testing.T) {
//	Convey("Given a CommandHandler", t, func() {
//		customers := new(mocks.PersistableCustomers)
//
//		repo := new(mocks.StartsRepositorySessions)
//		repo.On("StartSession", mock.AnythingOfType("*sql.Tx")).Return(customers)
//
//		db, _, err := sqlmock.New()
//		So(err, ShouldBeNil)
//
//		commandHandler := application.NewCommandHandler(repo, db)
//
//		Convey("When starting the repositry session fails", func() {
//			register, err := commands.NewRegister(
//				"64bcf656-da30-4f5a-b0b5-aead60965aa3",
//				"john@doe.com",
//				"John",
//				"Doe",
//			)
//			So(err, ShouldBeNil)
//
//			repo := new(mocks.StartsRepositorySessions)
//			db, _, err := sqlmock.New()
//			So(err, ShouldBeNil)
//			repo.On("StartSession").Return(nil, shared.ErrTechnical)
//			commandHandler = application.NewCommandHandler(repo, db)
//
//			err = commandHandler.Handle(register)
//
//			Convey("It should fail", func() {
//				So(err, ShouldBeError)
//				So(errors.Is(err, shared.ErrTechnical), ShouldBeTrue)
//			})
//		})
//
//		Convey("When committing the repositry session fails", func() {
//			register, err := commands.NewRegister(
//				"64bcf656-da30-4f5a-b0b5-aead60965aa3",
//				"john@doe.com",
//				"John",
//				"Doe",
//			)
//			So(err, ShouldBeNil)
//
//			customers.On("Register", mock.AnythingOfType("*domain.customer")).Return(nil).Once()
//			repo := new(mocks.StartsRepositorySessions)
//			repo.On("StartSession").Return(customers, nil)
//			db, _, err := sqlmock.New()
//			So(err, ShouldBeNil)
//			customers.On("Commit").Return(shared.ErrTechnical).Once()
//			commandHandler = application.NewCommandHandler(repo, db)
//
//			err = commandHandler.Handle(register)
//
//			Convey("It should fail", func() {
//				So(errors.Is(err, shared.ErrTechnical), ShouldBeTrue)
//				So(customers.AssertExpectations(t), ShouldBeTrue)
//			})
//		})
//
//		Convey("When rolling back the repositry session fails", func() {
//			register, err := commands.NewRegister(
//				"64bcf656-da30-4f5a-b0b5-aead60965aa3",
//				"john@doe.com",
//				"John",
//				"Doe",
//			)
//			So(err, ShouldBeNil)
//
//			customers.On("Register", mock.AnythingOfType("*domain.customer")).Return(shared.ErrDomainConstraintsViolation).Once()
//			repo := new(mocks.StartsRepositorySessions)
//			repo.On("StartSession").Return(customers, nil)
//			db, _, err := sqlmock.New()
//			So(err, ShouldBeNil)
//			customers.On("Rollback").Return(shared.ErrTechnical).Once()
//			commandHandler = application.NewCommandHandler(repo, db)
//
//			err = commandHandler.Handle(register)
//
//			Convey("It should fail", func() {
//				So(errors.Is(err, shared.ErrDomainConstraintsViolation), ShouldBeTrue)
//				So(customers.AssertExpectations(t), ShouldBeTrue)
//				fmt.Println(err)
//			})
//		})
//	})
//}
