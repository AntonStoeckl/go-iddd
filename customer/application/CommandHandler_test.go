package application_test

import (
	"errors"
	"fmt"
	"go-iddd/customer/application"
	"go-iddd/customer/application/mocks"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
	"golang.org/x/xerrors"
)

/*** Test factory method ***/

func TestNewCommandHandler(t *testing.T) {
	Convey("When a new CommandHandler is created", t, func() {
		repo := new(mocks.StartsRepositorySessions)
		commandHandler := application.NewCommandHandler(repo)

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
		repo := new(mocks.StartsRepositorySessions)
		customers := new(mocks.PersistableCustomersSession)
		repo.On("StartSession").Return(customers, nil)
		commandHandler := application.NewCommandHandler(repo)

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
					customers.On("Commit").Return(nil).Once()

					err := commandHandler.Handle(register)

					Convey("It should register and save a Customer", func() {
						So(err, ShouldBeNil)
						So(customers.AssertExpectations(t), ShouldBeTrue)
					})
				})

				Convey("And when saving the Customer fails", func() {
					expectedErr := errors.New("mocked error")
					customers.On("Register", mock.AnythingOfType("*domain.customer")).Return(expectedErr).Once()
					customers.On("Rollback").Return(nil).Once()

					err := commandHandler.Handle(register)

					Convey("It should fail", func() {
						So(xerrors.Is(err, expectedErr), ShouldBeTrue)
						So(customers.AssertExpectations(t), ShouldBeTrue)
					})
				})
			})
		})
	})
}

func TestCommandHandler_Handle_ConfirmEmailAddress(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		repo := new(mocks.StartsRepositorySessions)
		customers := new(mocks.PersistableCustomersSession)
		repo.On("StartSession").Return(customers, nil)
		commandHandler := application.NewCommandHandler(repo)

		Convey("And given a ConfirmEmailAddress command", func() {
			confirmEmailAddress, err := commands.NewConfirmEmailAddress(
				"64bcf656-da30-4f5a-b0b5-aead60965aa3",
				"john@doe.com",
				"secret_hash",
			)
			So(err, ShouldBeNil)

			conveyWhenTheCommandIsHandled(customers, confirmEmailAddress, commandHandler, t)
		})
	})
}

func TestCommandHandler_Handle_ChangeEmailAddress(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		repo := new(mocks.StartsRepositorySessions)
		customers := new(mocks.PersistableCustomersSession)
		repo.On("StartSession").Return(customers, nil)
		commandHandler := application.NewCommandHandler(repo)

		Convey("And given a ChangeEmailAddress command", func() {
			changeEmailAddress, err := commands.NewChangeEmailAddress(
				"64bcf656-da30-4f5a-b0b5-aead60965aa3",
				"john@doe.com",
			)
			So(err, ShouldBeNil)

			conveyWhenTheCommandIsHandled(customers, changeEmailAddress, commandHandler, t)
		})
	})
}

/*** Shared Convey for standard "modify" commands ***/

func conveyWhenTheCommandIsHandled(
	customers *mocks.PersistableCustomersSession,
	command shared.Command,
	commandHandler shared.CommandHandler,
	t *testing.T,
) {

	Convey("When the command is handled", func() {
		expectedErr := errors.New("mocked error")

		Convey("And when finding the Customer succeeds", func() {
			mockCustomer := new(mocks.Customer)
			customers.On("Of", command.AggregateID()).Return(mockCustomer, nil).Once()

			Convey("And when executing the command succeeds", func() {
				mockCustomer.On("Execute", command).Return(nil)

				Convey("And when saving the Customer succeeds", func() {
					customers.On("Persist", mockCustomer).Return(nil).Once()
					customers.On("Commit").Return(nil).Once()

					err := commandHandler.Handle(command)

					Convey("It should modify the Customer and save it", func() {
						So(err, ShouldBeNil)
						So(mockCustomer.AssertExpectations(t), ShouldBeTrue)
						So(customers.AssertExpectations(t), ShouldBeTrue)
					})
				})

				Convey("And when saving the Customer fails", func() {
					customers.On("Persist", mockCustomer).Return(expectedErr).Once()
					customers.On("Rollback").Return(nil).Once()

					err := commandHandler.Handle(command)

					Convey("It should fail", func() {
						So(xerrors.Is(err, expectedErr), ShouldBeTrue)
						So(customers.AssertExpectations(t), ShouldBeTrue)
					})
				})
			})

			Convey("And when executing the command fails", func() {
				mockCustomer.On("Execute", command).Return(expectedErr)
				customers.On("Rollback").Return(nil).Once()

				err := commandHandler.Handle(command)

				Convey("It should fail", func() {
					So(xerrors.Is(err, expectedErr), ShouldBeTrue)
					So(mockCustomer.AssertExpectations(t), ShouldBeTrue)
					So(customers.AssertExpectations(t), ShouldBeTrue)
				})
			})
		})

		Convey("And when finding the Customer fails", func() {
			customers.On("Of", command.AggregateID()).Return(nil, expectedErr).Once()
			customers.On("Rollback").Return(nil).Once()

			err := commandHandler.Handle(command)

			Convey("It should fail", func() {
				So(xerrors.Is(err, expectedErr), ShouldBeTrue)
				So(customers.AssertExpectations(t), ShouldBeTrue)
			})
		})
	})
}

/*** Test generic error cases ***/

func TestCommandHandler_Handle_WithInvalidCommand(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		repo := new(mocks.StartsRepositorySessions)
		customers := new(mocks.PersistableCustomersSession)
		repo.On("StartSession").Return(customers, nil)
		commandHandler := application.NewCommandHandler(repo)

		Convey("When a nil interface command is handled", func() {
			var nilInterfaceCommand shared.Command
			err := commandHandler.Handle(nilInterfaceCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
				fmt.Println(err)
			})
		})

		Convey("When a nil pointer command is handled", func() {
			var nilCommand *commands.ConfirmEmailAddress
			err := commandHandler.Handle(nilCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
				fmt.Println(err)
			})
		})

		Convey("When an empty command is handled", func() {
			emptyCommand := &commands.ConfirmEmailAddress{}
			err := commandHandler.Handle(emptyCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
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
				So(xerrors.Is(err, shared.ErrCommandIsUnknown), ShouldBeTrue)
				fmt.Println(err)
			})
		})
	})
}

func TestCommandHandler_Handle_WithSessionErrors(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		repo := new(mocks.StartsRepositorySessions)
		customers := new(mocks.PersistableCustomersSession)
		repo.On("StartSession").Return(customers, nil)
		commandHandler := application.NewCommandHandler(repo)

		Convey("When starting the repositry session fails", func() {
			register, err := commands.NewRegister(
				"64bcf656-da30-4f5a-b0b5-aead60965aa3",
				"john@doe.com",
				"John",
				"Doe",
			)
			So(err, ShouldBeNil)

			repo := new(mocks.StartsRepositorySessions)
			repo.On("StartSession").Return(nil, shared.ErrTechnical)
			commandHandler = application.NewCommandHandler(repo)

			err = commandHandler.Handle(register)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrTechnical), ShouldBeTrue)
				fmt.Println(err)
			})
		})

		Convey("When committing the repositry session fails", func() {
			register, err := commands.NewRegister(
				"64bcf656-da30-4f5a-b0b5-aead60965aa3",
				"john@doe.com",
				"John",
				"Doe",
			)
			So(err, ShouldBeNil)

			customers.On("Register", mock.AnythingOfType("*domain.customer")).Return(nil).Once()
			repo := new(mocks.StartsRepositorySessions)
			repo.On("StartSession").Return(customers, nil)
			expectedErr := errors.New("mocked error")
			customers.On("Commit").Return(expectedErr).Once()
			commandHandler = application.NewCommandHandler(repo)

			err = commandHandler.Handle(register)

			Convey("It should fail", func() {
				So(xerrors.Is(err, expectedErr), ShouldBeTrue)
				So(customers.AssertExpectations(t), ShouldBeTrue)
			})
		})

		Convey("When rolling back the repositry session fails", func() {
			register, err := commands.NewRegister(
				"64bcf656-da30-4f5a-b0b5-aead60965aa3",
				"john@doe.com",
				"John",
				"Doe",
			)
			So(err, ShouldBeNil)

			expectedErr := errors.New("mocked register error")
			customers.On("Register", mock.AnythingOfType("*domain.customer")).Return(expectedErr).Once()
			repo := new(mocks.StartsRepositorySessions)
			repo.On("StartSession").Return(customers, nil)
			expectedRollbackErr := errors.New("mocked rollback error")
			customers.On("Rollback").Return(expectedRollbackErr).Once()
			commandHandler = application.NewCommandHandler(repo)

			err = commandHandler.Handle(register)

			Convey("It should fail", func() {
				So(xerrors.Is(err, expectedErr), ShouldBeTrue)
				So(customers.AssertExpectations(t), ShouldBeTrue)
			})
		})
	})
}
