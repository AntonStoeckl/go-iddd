package customer_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	. "github.com/smartystreets/goconvey/convey"
)

func TestChangeName(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		var err error
		var recordedEvents es.DomainEvents

		customerID := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("kevin@ball.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.String())
		personName := values.RebuildPersonName("Kevin", "Ball")
		changedPersonName := values.RebuildPersonName("Latoya", "Ball")

		customerWasRegistered := events.BuildCustomerRegistered(
			customerID,
			emailAddress,
			confirmationHash,
			personName,
			1,
		)

		changeName, err := commands.BuildChangeCustomerName(
			customerID.String(),
			changedPersonName.GivenName(),
			changedPersonName.FamilyName(),
		)
		So(err, ShouldBeNil)

		Convey("\nSCENARIO 1: Change a Customer's name", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("When ChangeCustomerName", func() {
					recordedEvents, err = customer.ChangeName(eventStream, changeName)
					So(err, ShouldBeNil)

					Convey("Then CustomerNameChanged", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						nameChanged, ok := recordedEvents[0].(events.CustomerNameChanged)
						So(ok, ShouldBeTrue)
						So(nameChanged, ShouldNotBeNil)
						So(nameChanged.CustomerID().Equals(customerID), ShouldBeTrue)
						So(nameChanged.PersonName().Equals(changedPersonName), ShouldBeTrue)
						isError, reason := nameChanged.IndicatesAnError()
						So(isError, ShouldBeFalse)
						So(reason, ShouldBeBlank)
						So(nameChanged.Meta().StreamVersion(), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Try to change a Customer's name to the value he registered with", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("When ChangeCustomerName", func() {
					changeName, err = commands.BuildChangeCustomerName(
						customerID.String(),
						personName.GivenName(),
						personName.FamilyName(),
					)
					So(err, ShouldBeNil)

					recordedEvents, err = customer.ChangeName(eventStream, changeName)
					So(err, ShouldBeNil)

					Convey("Then no event", func() {
						So(recordedEvents, ShouldBeEmpty)
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Try to change a Customer's name to the value it was already changed to", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("and CustomerNameChanged", func() {
					nameChanged := events.BuildCustomerNameChanged(
						customerID,
						changedPersonName,
						2,
					)

					eventStream = append(eventStream, nameChanged)

					Convey("When ChangeCustomerName", func() {
						recordedEvents, err = customer.ChangeName(eventStream, changeName)
						So(err, ShouldBeNil)

						Convey("Then no event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 4: Try to change a Customer's name when the account was deleted", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("Given CustomerDeleted", func() {
					eventStream = append(
						eventStream,
						events.BuildCustomerDeleted(customerID, emailAddress, 2),
					)

					Convey("When ChangeCustomerName", func() {
						_, err := customer.ChangeName(eventStream, changeName)

						Convey("Then it should report an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
						})
					})
				})
			})
		})
	})
}
