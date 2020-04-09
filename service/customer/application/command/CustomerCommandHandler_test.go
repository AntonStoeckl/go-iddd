package command_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/application/command"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/customer/infrastructure/adapter/secondary/mocked"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

type commandHandlerTestArtifacts struct {
	emailAddress       string
	givenName          string
	familyName         string
	newEmailAddress    string
	newGivenName       string
	newFamilyName      string
	customerID         values.CustomerID
	confirmationHash   values.ConfirmationHash
	customerRegistered events.CustomerRegistered
}

func TestCustomerCommandHandler_ConcurrencyConflictHandling(t *testing.T) {
	customerEventStoreMock := new(mocked.ForStoringCustomerEvents)
	commandHandlerWithMock := command.NewCustomerCommandHandler(customerEventStoreMock)

	Convey("Prepare test artifacts", t, func() {
		var err error
		ca := buildArtifactsForCommandHandlerTest()

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
						err = commandHandlerWithMock.ConfirmCustomerEmailAddress(ca.customerID.String(), ca.confirmationHash.String())

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
						err = commandHandlerWithMock.ConfirmCustomerEmailAddress(ca.customerID.String(), ca.confirmationHash.String())

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrConcurrencyConflict), ShouldBeTrue)
						})
					})
				})
			})
		})
	})
}

func TestCustomerCommandHandler_TechnicalProblemsWithCustomerEventStore(t *testing.T) {
	customerEventStoreMock := new(mocked.ForStoringCustomerEvents)
	commandHandlerWithMock := command.NewCustomerCommandHandler(customerEventStoreMock)

	Convey("Prepare test artifacts", t, func() {
		var err error
		ca := buildArtifactsForCommandHandlerTest()

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
						err = commandHandlerWithMock.ConfirmCustomerEmailAddress(ca.customerID.String(), ca.confirmationHash.String())

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to change his email address", func() {
						err = commandHandlerWithMock.ChangeCustomerEmailAddress(ca.customerID.String(), ca.newEmailAddress)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to change his name", func() {
						err = commandHandlerWithMock.ChangeCustomerName(ca.customerID.String(), ca.givenName, ca.familyName)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to delete his account", func() {
						err = commandHandlerWithMock.DeleteCustomer(ca.customerID.String())

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
						err = commandHandlerWithMock.ConfirmCustomerEmailAddress(ca.customerID.String(), ca.confirmationHash.String())

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to change his email address", func() {
						err = commandHandlerWithMock.ChangeCustomerEmailAddress(ca.customerID.String(), ca.newEmailAddress)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to change his name", func() {
						err = commandHandlerWithMock.ChangeCustomerName(ca.customerID.String(), ca.givenName, ca.familyName)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to delete his account", func() {
						err = commandHandlerWithMock.DeleteCustomer(ca.customerID.String())

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})
				})
			})
		})
	})
}

func buildArtifactsForCommandHandlerTest() commandHandlerTestArtifacts {
	ca := commandHandlerTestArtifacts{}

	ca.emailAddress = "fiona@gallagher.net"
	ca.givenName = "Fiona"
	ca.familyName = "Galagher"
	ca.newEmailAddress = "fiona@pratt.net"
	ca.newGivenName = "Fiona"
	ca.newFamilyName = "Pratt"

	ca.customerID = values.GenerateCustomerID()
	ca.confirmationHash = values.GenerateConfirmationHash(ca.emailAddress)

	ca.customerRegistered = events.BuildCustomerRegistered(
		ca.customerID,
		values.RebuildEmailAddress(ca.emailAddress),
		ca.confirmationHash,
		values.RebuildPersonName(ca.givenName, ca.familyName),
		1,
	)

	return ca
}
