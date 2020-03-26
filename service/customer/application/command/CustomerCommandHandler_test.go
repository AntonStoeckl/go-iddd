package command_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/cmd"
	"github.com/AntonStoeckl/go-iddd/service/customer/application/command"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/secondary/mocked"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
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

		ca := buildArtifactsForCommandHandlerTest()

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
				err = commandHandler.RegisterCustomer(ca.registerCustomer)
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
				err = commandHandler.RegisterCustomer(ca.registerCustomer)
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
				err = commandHandler.RegisterCustomer(ca.registerCustomer)
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

		Convey("\nSCENARIO: Invalid command to delete a Customer's account", func() {
			Convey("Given a registered Customer", func() {
				err = commandHandler.RegisterCustomer(ca.registerCustomer)
				So(err, ShouldBeNil)

				Convey("When he tries to delete his account with an invalid command", func() {
					err = commandHandler.DeleteCustomer(commands.DeleteCustomer{})

					Convey("Then he should receive an error", func() {
						So(err, ShouldBeError)
						So(errors.Is(err, lib.ErrCommandIsInvalid), ShouldBeTrue)
					})
				})
			})
		})

		Convey("\nSCENARIO: Duplicate Customer ID", func() {
			Convey("Given a registered Customer", func() {
				err = commandHandler.RegisterCustomer(ca.registerCustomer)
				So(err, ShouldBeNil)

				Convey("And given he changed his email address", func() {
					err = commandHandler.ChangeCustomerEmailAddress(ca.changeCustomerEmailAddress)
					So(err, ShouldBeNil)

					Convey("When another prospective Customer tries to register and got a duplicate ID", func() {
						err = commandHandler.RegisterCustomer(ca.registerCustomer)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrDuplicate), ShouldBeTrue)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO: Technical problems with the CustomerEventStore", func() {
			Convey("Given a registered Customer", func() {
				Convey("and assuming the event stream can't be read", func() {
					customerEventStoreMock.
						On(
							"EventStreamFor",
							ca.customerID,
						).
						Return(nil, lib.ErrTechnical).
						Once()

					Convey("When he tries to confirm his email address", func() {
						err = commandHandlerWithMock.ConfirmCustomerEmailAddress(ca.confirmCustomerEmailAddress)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to change his email address", func() {
						err = commandHandlerWithMock.ChangeCustomerEmailAddress(ca.changeCustomerEmailAddress)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to change his name", func() {
						err = commandHandlerWithMock.ChangeCustomerName(ca.changeCustomerName)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to delete his account", func() {
						err = commandHandlerWithMock.DeleteCustomer(ca.deleteCustomer)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})
				})

				Convey("and assuming the recorded events can't be stored", func() {
					customerEventStoreMock.
						On("EventStreamFor", ca.customerID).
						Return(es.DomainEvents{ca.customerRegistered}, nil).
						Once()

					customerEventStoreMock.
						On(
							"Add",
							mock.AnythingOfType("es.DomainEvents"),
							ca.customerID,
						).
						Return(lib.ErrTechnical).
						Once()

					Convey("When he tries to confirm his email address", func() {
						err = commandHandlerWithMock.ConfirmCustomerEmailAddress(ca.confirmCustomerEmailAddress)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to change his email address", func() {
						err = commandHandlerWithMock.ChangeCustomerEmailAddress(ca.changeCustomerEmailAddress)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to change his name", func() {
						err = commandHandlerWithMock.ChangeCustomerName(ca.changeCustomerName)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to delete his account", func() {
						err = commandHandlerWithMock.DeleteCustomer(ca.deleteCustomer)

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
					On("EventStreamFor", ca.customerID).
					Return(es.DomainEvents{ca.customerRegistered}, nil).
					Times(12)

				Convey("and assuming a concurrency conflict happens once", func() {
					customerEventStoreMock.
						On(
							"Add",
							mock.AnythingOfType("es.DomainEvents"),
							ca.customerID,
						).
						Return(lib.ErrConcurrencyConflict).
						Once()

					customerEventStoreMock.
						On(
							"Add",
							mock.AnythingOfType("es.DomainEvents"),
							ca.customerID,
						).
						Return(nil).
						Once()

					Convey("When he tries to confirm his email address", func() {
						err = commandHandlerWithMock.ConfirmCustomerEmailAddress(ca.confirmCustomerEmailAddress)

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
							ca.customerID,
						).
						Return(lib.ErrConcurrencyConflict).
						Times(10)

					Convey("When he tries to confirm his email address", func() {
						err = commandHandlerWithMock.ConfirmCustomerEmailAddress(ca.confirmCustomerEmailAddress)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrConcurrencyConflict), ShouldBeTrue)
						})
					})
				})
			})
		})

		Reset(func() {
			err = commandHandler.DeleteCustomer(ca.deleteCustomer)
			if !errors.Is(err, lib.ErrNotFound) {
				So(err, ShouldBeNil)
			}

			err = diContainer.GetCustomerEventStore().Delete(ca.customerID)
			So(err, ShouldBeNil)
		})
	})
}

type commandHandlerTestArtifacts struct {
	customerID                  values.CustomerID
	registerCustomer            commands.RegisterCustomer
	confirmCustomerEmailAddress commands.ConfirmCustomerEmailAddress
	changeCustomerEmailAddress  commands.ChangeCustomerEmailAddress
	changeCustomerName          commands.ChangeCustomerName
	deleteCustomer              commands.DeleteCustomer
	customerRegistered          events.CustomerRegistered
}

func buildArtifactsForCommandHandlerTest() commandHandlerTestArtifacts {
	var err error

	ca := commandHandlerTestArtifacts{}

	ca.registerCustomer, err = commands.BuildRegisterCustomer(
		"john@doe.com",
		"John",
		"Doe",
	)
	So(err, ShouldBeNil)

	ca.customerID = ca.registerCustomer.CustomerID()

	ca.confirmCustomerEmailAddress, err = commands.BuildConfirmCustomerEmailAddress(
		ca.customerID.ID(),
		ca.registerCustomer.ConfirmationHash().Hash(),
	)
	So(err, ShouldBeNil)

	ca.changeCustomerEmailAddress, err = commands.BuildChangeCustomerEmailAddress(
		ca.customerID.ID(),
		"john+changed@doe.com",
	)
	So(err, ShouldBeNil)

	ca.changeCustomerName, err = commands.BuildChangeCustomerName(
		ca.customerID.ID(),
		"James",
		"Dope",
	)
	So(err, ShouldBeNil)

	ca.deleteCustomer, err = commands.BuildDeleteCustomer(ca.customerID.ID())
	So(err, ShouldBeNil)

	ca.customerRegistered = events.CustomerWasRegistered(
		ca.customerID,
		ca.registerCustomer.EmailAddress(),
		ca.registerCustomer.ConfirmationHash(),
		ca.registerCustomer.PersonName(),
		uint(1),
	)

	return ca
}
