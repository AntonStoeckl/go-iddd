package command_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/application/command"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
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

func TestCustomerCommandHandler_TechnicalProblemsWithCustomerEventStore(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		var err error
		ca := buildArtifactsForCommandHandlerTest()

		Convey("\nSCENARIO: Technical problems with the CustomerEventStore", func() {
			Convey("Given a registered Customer", func() {
				Convey("and assuming the event stream can't be read", func() {
					commandHandler := command.NewCustomerCommandHandler(
						func(id values.CustomerID) (es.EventStream, error) {
							return nil, lib.ErrTechnical
						},
						func(recordedEvents es.RecordedEvents, id values.CustomerID) error {
							return nil
						},
						func(recordedEvents es.RecordedEvents, id values.CustomerID) error {
							return nil
						},
						lib.RetryOnConcurrencyConflict,
					)

					Convey("When he tries to confirm his email address", func() {
						err = commandHandler.ConfirmCustomerEmailAddress(ca.customerID.String(), ca.confirmationHash.String())

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to change his email address", func() {
						err = commandHandler.ChangeCustomerEmailAddress(ca.customerID.String(), ca.newEmailAddress)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to change his name", func() {
						err = commandHandler.ChangeCustomerName(ca.customerID.String(), ca.givenName, ca.familyName)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to delete his account", func() {
						err = commandHandler.DeleteCustomer(ca.customerID.String())

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})
				})

				Convey("and assuming the recorded events can't be stored", func() {
					commandHandler := command.NewCustomerCommandHandler(
						func(id values.CustomerID) (es.EventStream, error) {
							return es.EventStream{ca.customerRegistered}, nil
						},
						func(recordedEvents es.RecordedEvents, id values.CustomerID) error {
							return nil
						},
						func(recordedEvents es.RecordedEvents, id values.CustomerID) error {
							return lib.ErrTechnical
						},
						lib.RetryOnConcurrencyConflict,
					)

					Convey("When he tries to confirm his email address", func() {
						err = commandHandler.ConfirmCustomerEmailAddress(ca.customerID.String(), ca.confirmationHash.String())

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to change his email address", func() {
						err = commandHandler.ChangeCustomerEmailAddress(ca.customerID.String(), ca.newEmailAddress)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to change his name", func() {
						err = commandHandler.ChangeCustomerName(ca.customerID.String(), ca.givenName, ca.familyName)

						Convey("Then he should receive an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrTechnical), ShouldBeTrue)
						})
					})

					Convey("When he tries to delete his account", func() {
						err = commandHandler.DeleteCustomer(ca.customerID.String())

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
