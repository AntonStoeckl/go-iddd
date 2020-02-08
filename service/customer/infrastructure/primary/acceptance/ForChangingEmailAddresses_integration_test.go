package acceptance_test

import (
	"go-iddd/service/customer/application"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/customer/infrastructure"
	"go-iddd/service/customer/infrastructure/secondary/forstoringcustomerevents/mocked"
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func Test_ChangeEmailAddress(t *testing.T) {
	Convey("Setup", t, func() {
		diContainer, err := infrastructure.SetUpDIContainer()
		So(err, ShouldBeNil)
		commandHandler := diContainer.GetCustomerCommandHandler()
		customerEventStore := diContainer.GetCustomerEventStore()

		newEmailAddress := "john+changed@doe.com"

		Convey("Given a registered Customer", func() {
			register, err := commands.NewRegister(
				"john@doe.com",
				"John",
				"Doe",
			)
			So(err, ShouldBeNil)

			err = commandHandler.Register(register)
			So(err, ShouldBeNil)

			Convey("When a Customer's emailAddress is changed", func() {
				changeEmailAddress, err := commands.NewChangeEmailAddress(
					register.CustomerID().ID(),
					newEmailAddress,
				)
				So(err, ShouldBeNil)

				err = commandHandler.ChangeEmailAddress(changeEmailAddress)

				Convey("It should succeed", func() {
					So(err, ShouldBeNil)

					eventStream, err := customerEventStore.EventStreamFor(register.CustomerID())
					So(err, ShouldBeNil)
					So(eventStream, ShouldHaveLength, 2)
					So(eventStream[0], ShouldHaveSameTypeAs, events.Registered{})
					So(eventStream[1], ShouldHaveSameTypeAs, events.EmailAddressChanged{})

					Convey("And when a Customer's emailAddress is changed again to the same emailAddress", func() {
						changeEmailAddress, err := commands.NewChangeEmailAddress(
							register.CustomerID().ID(),
							newEmailAddress,
						)
						So(err, ShouldBeNil)

						err = commandHandler.ChangeEmailAddress(changeEmailAddress)

						Convey("It should succeed", func() {
							So(err, ShouldBeNil)

							eventStream, err := customerEventStore.EventStreamFor(register.CustomerID())
							So(err, ShouldBeNil)
							So(eventStream, ShouldHaveLength, 2)
							So(eventStream[0], ShouldHaveSameTypeAs, events.Registered{})
							So(eventStream[1], ShouldHaveSameTypeAs, events.EmailAddressChanged{})
						})
					})
				})
			})

			Convey("When a Customer's emailAddress is changed with an invalid command", func() {
				err := commandHandler.ChangeEmailAddress(commands.ChangeEmailAddress{})

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
				})
			})

			err = diContainer.GetCustomerEventStore().Delete(register.CustomerID())
			So(err, ShouldBeNil)
		})

		Convey("Given an unregistered Customer", func() {
			Convey("When a Customer's emailAddress is changed", func() {
				changeEmailAddress, err := commands.NewChangeEmailAddress(
					values.GenerateCustomerID().ID(),
					newEmailAddress,
				)
				So(err, ShouldBeNil)

				err = commandHandler.ChangeEmailAddress(changeEmailAddress)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
				})
			})
		})
	})
}

func Test_ChangeEmailAddress_WithErrorFromCustomers(t *testing.T) {
	Convey("Setup", t, func() {
		customerEventStore := new(mocked.ForStoringCustomerEvents)

		commandHandler := application.NewCommandHandler(customerEventStore)

		customerID := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("john@doe.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("John", "Doe")

		eventStream := es.DomainEvents{
			events.ItWasRegistered(customerID, emailAddress, confirmationHash, personName, uint(1)),
		}

		changeEmailAddress, err := commands.NewChangeEmailAddress(
			customerID.ID(),
			"john+changed@doe.com",
		)
		So(err, ShouldBeNil)

		Convey("Given the Customer can't be persisted because of a technical error", func() {
			customerEventStore.On("EventStreamFor", mock.Anything).Return(eventStream, nil).Once()
			customerEventStore.On("Add", mock.Anything, mock.Anything, mock.Anything).Return(lib.ErrTechnical).Once()

			Convey("When a Customer's emailAddress is confirmed", func() {
				err := commandHandler.ChangeEmailAddress(changeEmailAddress)

				Convey("It should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
				})
			})
		})
	})
}

func Test_ChangeEmailAddress_RetriesWhenItHasConcurrencyConflicts(t *testing.T) {
	Convey("Given a CommandHandler", t, func() {
		customerEventStore := new(mocked.ForStoringCustomerEvents)

		commandHandler := application.NewCommandHandler(customerEventStore)

		customerID := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("john@doe.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("John", "Doe")
		eventStream := es.DomainEvents{
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
						customerEventStore.On("EventStreamFor", mock.Anything).Return(eventStream, nil).Twice()

						// fist attempt runs into a concurrency conflict
						customerEventStore.
							On("Add", mock.Anything, mock.Anything, mock.Anything).
							Return(lib.ErrConcurrencyConflict).
							Once()

						// second attempt works
						customerEventStore.
							On("Add", mock.Anything, mock.Anything, mock.Anything).
							Return(nil).
							Once()

						err := commandHandler.ChangeEmailAddress(changeEmailAddress)

						Convey("It should modify the Customer and save it", func() {
							So(err, ShouldBeNil)
						})
					})

					Convey("And when saving the Customer has too many concurrency conflicts", func() {
						// should be called 10 times due to retries
						customerEventStore.
							On("EventStreamFor", mock.Anything).
							Return(eventStream, nil).
							Times(10)

						// all attempts run into a concurrency conflict
						customerEventStore.
							On("Add", mock.Anything, mock.Anything, mock.Anything).
							Return(lib.ErrConcurrencyConflict).
							Times(10)

						err := commandHandler.ChangeEmailAddress(changeEmailAddress)

						Convey("It should fail", func() {
							So(err, ShouldNotBeNil)
							So(errors.Is(err, lib.ErrConcurrencyConflict), ShouldBeTrue)
						})
					})
				})
			})
		})
	})
}
