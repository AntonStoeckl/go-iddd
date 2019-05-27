package application_test

import (
	"errors"
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
		mockCustomers := new(mocks.Customers)
		commandHandler := application.NewCommandHandler(mockCustomers)

		Convey("It should succeed", func() {
			So(commandHandler, ShouldNotBeNil)
			So(commandHandler, ShouldImplement, (*shared.CommandHandler)(nil))
			So(commandHandler, ShouldHaveSameTypeAs, (*application.CommandHandler)(nil))
		})
	})
}

/*** Test business cases ***/

func TestHandleRegister(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		mockCustomers := new(mocks.Customers)
		commandHandler := application.NewCommandHandler(mockCustomers)

		Convey("And given a Register command", func() {
			id := "64bcf656-da30-4f5a-b0b5-aead60965aa3"
			emailAddress := "john@doe.com"
			givenName := "John"
			familyName := "Doe"

			register, err := commands.NewRegister(id, emailAddress, givenName, familyName)
			So(err, ShouldBeNil)

			Convey("When the command is handled", func() {
				Convey("And when saving the Customer succeeds", func() {
					mockCustomers.On("Register", mock.AnythingOfType("*domain.customer")).Return(nil).Once()
					err := commandHandler.Handle(register)

					Convey("It should register and save a Customer", func() {
						So(err, ShouldBeNil)
						So(mockCustomers.AssertExpectations(t), ShouldBeTrue)
					})
				})

				Convey("And when saving the Customer fails", func() {
					expectedErr := errors.New("mocked error")
					mockCustomers.On("Register", mock.AnythingOfType("*domain.customer")).Return(expectedErr).Once()
					err := commandHandler.Handle(register)

					Convey("It should fail", func() {
						So(err, ShouldBeError, expectedErr)
						So(mockCustomers.AssertExpectations(t), ShouldBeTrue)
					})
				})
			})
		})
	})
}

func TestHandleConfirmEmailAddress(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		mockCustomers := new(mocks.Customers)
		commandHandler := application.NewCommandHandler(mockCustomers)

		Convey("And given a ConfirmEmailAddress command", func() {
			id := "64bcf656-da30-4f5a-b0b5-aead60965aa3"
			emailAddress := "john@doe.com"
			confirmationHash := "secret_hash"

			confirmEmailAddress, err := commands.NewConfirmEmailAddress(id, emailAddress, confirmationHash)
			So(err, ShouldBeNil)

			Convey("When the command is handled", func() {
				expectedErr := errors.New("mocked error")

				Convey("And when finding the Customer succeeds", func() {
					mockCustomer := new(mocks.Customer)
					mockCustomers.On("Of", confirmEmailAddress.ID()).Return(mockCustomer, nil).Once()

					Convey("And when applying confirmEmailAddress succeeds", func() {
						mockCustomer.On("Apply", confirmEmailAddress).Return(nil)

						Convey("And when saving the Customer succeeds", func() {
							mockCustomers.On("Save", mockCustomer).Return(nil).Once()
							err := commandHandler.Handle(confirmEmailAddress)

							Convey("It should confirmEmailAddress of a Customer and save it", func() {
								So(err, ShouldBeNil)
								So(mockCustomer.AssertExpectations(t), ShouldBeTrue)
								So(mockCustomers.AssertExpectations(t), ShouldBeTrue)
							})
						})

						Convey("And when saving the Customer fails", func() {
							mockCustomers.On("Save", mockCustomer).Return(expectedErr).Once()
							err := commandHandler.Handle(confirmEmailAddress)

							Convey("It should fail", func() {
								So(err, ShouldBeError, expectedErr)
								So(mockCustomers.AssertExpectations(t), ShouldBeTrue)
							})
						})
					})

					Convey("And when applying confirmEmailAddress fails", func() {
						mockCustomer.On("Apply", confirmEmailAddress).Return(expectedErr)
						err := commandHandler.Handle(confirmEmailAddress)

						Convey("It should fail", func() {
							So(err, ShouldBeError, expectedErr)
							So(mockCustomer.AssertExpectations(t), ShouldBeTrue)
							So(mockCustomers.AssertExpectations(t), ShouldBeTrue)
						})
					})
				})

				Convey("And when finding the Customer fails", func() {
					mockCustomers.On("Of", confirmEmailAddress.ID()).Return(nil, expectedErr).Once()
					err := commandHandler.Handle(confirmEmailAddress)

					Convey("It should fail", func() {
						So(err, ShouldBeError, expectedErr)
						So(mockCustomers.AssertExpectations(t), ShouldBeTrue)
					})
				})
			})
		})
	})
}

/*** Test handling invalid commands ***/

func TestHandleInvalidCommand(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		mockCustomers := new(mocks.Customers)
		commandHandler := application.NewCommandHandler(mockCustomers)
		So(commandHandler, ShouldImplement, (*shared.CommandHandler)(nil))

		Convey("When a nil interface command is handled", func() {
			var nilInterfaceCommand shared.Command
			err := commandHandler.Handle(nilInterfaceCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
			})
		})

		Convey("When a nil pointer command is handled", func() {
			var nilCommand *commands.ConfirmEmailAddress
			err := commandHandler.Handle(nilCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
			})
		})

		Convey("When an empty command is handled", func() {
			emptyCommand := &commands.ConfirmEmailAddress{}
			err := commandHandler.Handle(emptyCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrCommandIsInvalid), ShouldBeTrue)
			})
		})

		Convey("When an unknown command is handled", func() {
			unknownCommand := &unknownCommand{}
			err := commandHandler.Handle(unknownCommand)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrCommandCanNotBeHandled), ShouldBeTrue)
			})
		})
	})
}

/*** Test Helpers ***/

type unknownCommand struct{}

func (c *unknownCommand) AggregateIdentifier() shared.AggregateIdentifier {
	return values.GenerateID()
}

func (c *unknownCommand) CommandName() string {
	return "unknown"
}
