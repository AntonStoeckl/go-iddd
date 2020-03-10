package application_test

import (
	"go-iddd/service/customer/application"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/customer/infrastructure/secondary/forstoringcustomerevents/mocked"
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestCommandHandler(t *testing.T) {
	Convey("Setup", t, func() {
		customerEventStore := new(mocked.ForStoringCustomerEvents)
		commandHandler := application.NewCommandHandler(customerEventStore)

		Convey("\nSCENARIO 1: Invalid Commands", func() {
			Convey("When a Customer is registered with an invalid Command", func() {
				err := commandHandler.RegisterCustomer(commands.RegisterCustomer{})

				Convey("Then it should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
				})
			})

			Convey("When a Customer's emailAddress is confirmed with an invalid command", func() {
				err := commandHandler.ConfirmCustomerEmailAddress(commands.ConfirmEmailAddress{})

				Convey("Then it should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
				})
			})

			Convey("When a Customer's emailAddress is changed with an invalid command", func() {
				err := commandHandler.ChangeCustomerEmailAddress(commands.ChangeEmailAddress{})

				Convey("Then it should fail", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
				})
			})
		})

		Convey("\nSCENARIO 2: Duplicate Customer ID", func() {
			Convey("Given a registered Customer", func() {
				register, err := commands.BuildRegisterCustomer(
					"john@doe.com",
					"John",
					"Doe",
				)
				So(err, ShouldBeNil)

				customerEventStore.
					On(
						"CreateStreamFrom",
						mock.AnythingOfType("es.DomainEvents"),
						register.CustomerID(),
					).
					Return(lib.ErrDuplicate).
					Once()

				Convey("When he is registered again with duplicate ID", func() {
					err = commandHandler.RegisterCustomer(register)

					Convey("Then it should fail", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrDuplicate), ShouldBeTrue)
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Customer does not exist", func() {
			Convey("Given an unregistered Customer", func() {
				confirmEmailAddress, err := commands.NewConfirmEmailAddress(
					values.GenerateCustomerID().ID(),
					values.GenerateConfirmationHash("john@doe.com").Hash(),
				)
				So(err, ShouldBeNil)

				customerEventStore.
					On("EventStreamFor", confirmEmailAddress.CustomerID()).
					Return(es.DomainEvents{}, lib.ErrNotFound).
					Once()

				Convey("When his emailAddress is confirmed", func() {
					err = commandHandler.ConfirmCustomerEmailAddress(confirmEmailAddress)

					Convey("Then it should fail", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
					})
				})

				changeEmailAddress, err := commands.NewChangeEmailAddress(
					values.GenerateCustomerID().ID(),
					"john@doe.com",
				)
				So(err, ShouldBeNil)

				customerEventStore.
					On("EventStreamFor", changeEmailAddress.CustomerID()).
					Return(es.DomainEvents{}, lib.ErrNotFound).
					Once()

				Convey("When his emailAddress is changed", func() {
					err = commandHandler.ChangeCustomerEmailAddress(changeEmailAddress)

					Convey("Then it should fail", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
					})
				})
			})
		})

		Convey("\nSCENARIO 4: Technical problems with the CustomerEventStore", func() {
			Convey("Given a registered Customer", func() {
				registered := events.CustomerWasRegistered(
					values.GenerateCustomerID(),
					values.RebuildEmailAddress("john@doe.com"),
					values.RebuildConfirmationHash("john@doe.com"),
					values.RebuildPersonName("John", "Doe"),
					uint(1),
				)

				customerEventStore.
					On("EventStreamFor", registered.CustomerID()).
					Return(es.DomainEvents{registered}, nil).
					Once()

				Convey("and assuming the recorded events can't be stored", func() {
					customerEventStore.
						On(
							"Add",
							mock.AnythingOfType("es.DomainEvents"),
							registered.CustomerID(),
						).
						Return(lib.ErrTechnical).
						Once()

					Convey("When trying to confirm his emailAddress", func() {
						confirmEmailAddress, err := commands.NewConfirmEmailAddress(
							registered.CustomerID().ID(),
							registered.ConfirmationHash().Hash(),
						)
						So(err, ShouldBeNil)

						err = commandHandler.ConfirmCustomerEmailAddress(confirmEmailAddress)

						Convey("Then it should fail", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When trying to change his emailAddress", func() {
						changeEmailAddress, err := commands.NewChangeEmailAddress(
							registered.CustomerID().ID(),
							registered.EmailAddress().EmailAddress(),
						)
						So(err, ShouldBeNil)

						err = commandHandler.ChangeCustomerEmailAddress(changeEmailAddress)

						Convey("Then it should fail", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})
				})
			})
		})
	})
}
