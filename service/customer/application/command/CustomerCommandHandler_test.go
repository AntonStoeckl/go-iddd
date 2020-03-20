package command_test

import (
	"go-iddd/service/cmd"
	"go-iddd/service/customer/application/command"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/infrastructure/secondary/mocked"
	"go-iddd/service/lib"
	"go-iddd/service/lib/es"
	"testing"

	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestCustomerCommandHandler(t *testing.T) {
	diContainer, err := cmd.Bootstrap()
	if err != nil {
		panic(err)
	}

	commandHandler := diContainer.GetCustomerCommandHandler()

	customerEventStoreMock := new(mocked.ForStoringCustomerEvents)
	commandHandlerWithMock := command.NewCustomerCommandHandler(customerEventStoreMock)

	Convey("Prepare test artifacts", t, func() {
		var err error

		registerCustomer, err := commands.BuildRegisterCustomer(
			"john@doe.com",
			"John",
			"Doe",
		)
		So(err, ShouldBeNil)

		customerID := registerCustomer.CustomerID()

		confirmCustomerEmailAddress, err := commands.BuildConfirmCustomerEmailAddress(
			customerID.ID(),
			registerCustomer.ConfirmationHash().Hash(),
		)
		So(err, ShouldBeNil)

		changeCustomerEmailAddress, err := commands.BuildChangeCustomerEmailAddress(
			customerID.ID(),
			"john+changed@doe.com",
		)
		So(err, ShouldBeNil)

		changeCustomerName, err := commands.BuildChangeCustomerName(
			customerID.ID(),
			"James",
			"Dope",
		)
		So(err, ShouldBeNil)

		registered := events.CustomerWasRegistered(
			customerID,
			registerCustomer.EmailAddress(),
			registerCustomer.ConfirmationHash(),
			registerCustomer.PersonName(),
			uint(1),
		)

		Convey("\nSCENARIO: Invalid command to register a prospective Customer", func() {
			Convey("When a Customer registers with an invalid Command", func() {
				err = commandHandler.RegisterCustomer(commands.RegisterCustomer{})

				Convey("Then he should receive an error", func() {
					So(err, ShouldBeError)
					So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
				})
			})
		})

		Convey("\nSCENARIO: Invalid command to confirm a Customer's email address", func() {
			Convey("Given a registered Customer", func() {
				err = commandHandler.RegisterCustomer(registerCustomer)
				So(err, ShouldBeNil)

				Convey("When he tries to confirm his email address with an invalid command", func() {
					err = commandHandler.ConfirmCustomerEmailAddress(commands.ConfirmCustomerEmailAddress{})

					Convey("Then he should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
					})
				})
			})
		})

		Convey("\nSCENARIO: Invalid command to change a Customer's email address", func() {
			Convey("Given a registered Customer", func() {
				err = commandHandler.RegisterCustomer(registerCustomer)
				So(err, ShouldBeNil)

				Convey("When he tries to change his email address with an invalid command", func() {
					err = commandHandler.ChangeCustomerEmailAddress(commands.ChangeCustomerEmailAddress{})

					Convey("Then he should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
					})
				})
			})
		})

		Convey("\nSCENARIO: Invalid command to change a Customer's name", func() {
			Convey("Given a registered Customer", func() {
				err = commandHandler.RegisterCustomer(registerCustomer)
				So(err, ShouldBeNil)

				Convey("When he tries to change his name with an invalid command", func() {
					err = commandHandler.ChangeCustomerName(commands.ChangeCustomerName{})

					Convey("Then he should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
					})
				})
			})
		})

		Convey("\nSCENARIO: Duplicate Customer ID", func() {
			Convey("Given a registered Customer", func() {
				err = commandHandler.RegisterCustomer(registerCustomer)
				So(err, ShouldBeNil)

				Convey("When another prospective Customer tries to register and got a duplicate ID", func() {
					err = commandHandler.RegisterCustomer(registerCustomer)

					Convey("Then he should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrDuplicate), ShouldBeTrue)
					})
				})
			})
		})

		Convey("\nSCENARIO: Technical problems with the CustomerEventStore", func() {
			Convey("Given a registered Customer", func() {
				customerEventStoreMock.
					On("EventStreamFor", customerID).
					Return(es.DomainEvents{registered}, nil).
					Once()

				Convey("and assuming the recorded events can't be stored", func() {
					customerEventStoreMock.
						On(
							"Add",
							mock.AnythingOfType("es.DomainEvents"),
							customerID,
						).
						Return(lib.ErrTechnical).
						Once()

					Convey("When he tries to confirm his email address", func() {
						err = commandHandlerWithMock.ConfirmCustomerEmailAddress(confirmCustomerEmailAddress)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to change his email address", func() {
						err = commandHandlerWithMock.ChangeCustomerEmailAddress(changeCustomerEmailAddress)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to change his name", func() {
						err = commandHandlerWithMock.ChangeCustomerName(changeCustomerName)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: Concurrency conflict in CustomerEventStore", func() {
			Convey("Given a registered Customer", func() {
				customerEventStoreMock.
					On("EventStreamFor", customerID).
					Return(es.DomainEvents{registered}, nil).
					Times(12)

				Convey("and assuming a concurrency conflict happens once", func() {
					customerEventStoreMock.
						On(
							"Add",
							mock.AnythingOfType("es.DomainEvents"),
							customerID,
						).
						Return(lib.ErrConcurrencyConflict).
						Once()

					customerEventStoreMock.
						On(
							"Add",
							mock.AnythingOfType("es.DomainEvents"),
							customerID,
						).
						Return(nil).
						Once()

					Convey("When he tries to confirm his email address", func() {
						err = commandHandlerWithMock.ConfirmCustomerEmailAddress(confirmCustomerEmailAddress)

						Convey("Then it should succeed after retry", func() {
							So(err, ShouldBeNil)
						})
					})
				})

				Convey("and assuming a concurrency conflict happens 10 times", func() {
					customerEventStoreMock.
						On(
							"Add",
							mock.AnythingOfType("es.DomainEvents"),
							customerID,
						).
						Return(lib.ErrConcurrencyConflict).
						Times(10)

					Convey("When he tries to confirm his email address", func() {
						err = commandHandlerWithMock.ConfirmCustomerEmailAddress(confirmCustomerEmailAddress)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrConcurrencyConflict), ShouldBeTrue)
						})
					})
				})
			})
		})

		Reset(func() {
			err := diContainer.GetCustomerEventStore().Delete(customerID)
			So(err, ShouldBeNil)
		})
	})
}
