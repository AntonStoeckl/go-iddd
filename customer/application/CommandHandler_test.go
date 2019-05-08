package application

import (
	"errors"
	"go-iddd/customer/application/mocks"
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/valueobjects"
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHandleRegister(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		mockCustomers := new(mocks.Customers)
		commandHandler := NewCommandHandler(mockCustomers)

		Convey("And given a Register command", func() {
			id := valueobjects.GenerateID()
			emailAddress, err := valueobjects.NewEmailAddress("foo@bar.com")
			So(err, ShouldBeNil)
			personName, err := valueobjects.NewPersonName("John", "Doe")
			So(err, ShouldBeNil)

			register, err := domain.NewRegister(id, emailAddress, personName)
			So(err, ShouldBeNil)

			Convey("When the command is handled", func() {
				mockCustomer := new(mocks.Customer)
				mockCustomers.On("New").Return(mockCustomer)
				expectedErr := errors.New("mocked error")

				Convey("And when applying register succeeds", func() {
					mockCustomer.On("Apply", register).Return(nil)

					Convey("And when saving the Customer succeeds", func() {
						mockCustomers.On("Save", mockCustomer).Return(nil).Once()

						err := commandHandler.Handle(register)

						Convey("It should register and save a Customer", func() {
							So(err, ShouldBeNil)
							So(mockCustomer.AssertExpectations(t), ShouldBeTrue)
							So(mockCustomers.AssertExpectations(t), ShouldBeTrue)
						})
					})

					Convey("And when saving the Customer fails", func() {
						expectedErr := errors.New("mocked error")
						mockCustomers.On("Save", mockCustomer).Return(expectedErr).Once()

						err := commandHandler.Handle(register)

						Convey("It should fail", func() {
							So(err, ShouldBeError, expectedErr)
							So(mockCustomers.AssertExpectations(t), ShouldBeTrue)
						})
					})
				})

				Convey("And when applying register fails", func() {
					mockCustomer.On("Apply", register).Return(expectedErr)

					err := commandHandler.Handle(register)

					Convey("It should fail", func() {
						So(err, ShouldBeError, expectedErr)
						So(mockCustomer.AssertExpectations(t), ShouldBeTrue)
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
		commandHandler := NewCommandHandler(mockCustomers)

		Convey("And given a ConfirmEmailAddress command", func() {
			id := valueobjects.GenerateID()
			emailAddress, err := valueobjects.NewEmailAddress("foo@bar.com")
			So(err, ShouldBeNil)
			confirmationHash := valueobjects.GenerateConfirmationHash(emailAddress.EmailAddress())

			confirmEmailAddress, err := domain.NewConfirmEmailAddress(id, emailAddress, confirmationHash)
			So(err, ShouldBeNil)

			Convey("When the command is handled", func() {
				expectedErr := errors.New("mocked error")

				Convey("And when finding the Customer succeeds", func() {
					mockCustomer := new(mocks.Customer)
					mockCustomers.On("FindBy", id).Return(mockCustomer, nil).Once()

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
					mockCustomers.On("FindBy", id).Return(nil, expectedErr).Once()
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

func TestHandleInvalidCommand(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		mockCustomers := new(mocks.Customers)
		commandHandler := NewCommandHandler(mockCustomers)
		So(commandHandler, ShouldImplement, (*shared.CommandHandler)(nil))

		Convey("And given an invalid Command", func() {
			var invalidCommand shared.Command

			Convey("When the command is handled", func() {
				err := commandHandler.Handle(invalidCommand)

				Convey("It should fail", func() {
					So(err, ShouldBeError, "commandHandler - nil command handled")
				})
			})
		})
	})
}

func TestHandleUnknownCommand(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		mockCustomers := new(mocks.Customers)
		commandHandler := NewCommandHandler(mockCustomers)
		So(commandHandler, ShouldImplement, (*shared.CommandHandler)(nil))

		Convey("And given an unknown Command", func() {
			unknownCommand := &unknownCommand{}

			Convey("When the command is handled", func() {
				err := commandHandler.Handle(unknownCommand)

				Convey("It should fail", func() {
					So(err, ShouldBeError, "commandHandler - unknown command handled")
				})
			})
		})
	})
}

type unknownCommand struct{}

func (c *unknownCommand) AggregateIdentifier() shared.AggregateIdentifier {
	return valueobjects.GenerateID()
}

func (c *unknownCommand) CommandName() string {
	return "unknown"
}
