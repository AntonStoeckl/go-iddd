package customer_test

import (
	"go-iddd/service/customer/domain/customer"
	"go-iddd/service/customer/domain/customer/commands"
	"go-iddd/service/customer/domain/customer/events"
	"go-iddd/service/customer/domain/customer/values"
	"go-iddd/service/lib/es"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChangeName(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		var err error

		customerID := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("kevin@ball.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("Kevin", "Ball")
		changedPersonName := values.RebuildPersonName("Latoya", "Ball")

		customerWasRegistered := events.CustomerWasRegistered(
			customerID,
			emailAddress,
			confirmationHash,
			personName,
			1,
		)

		changeName, err := commands.BuildChangeCustomerName(
			customerID.ID(),
			changedPersonName.GivenName(),
			changedPersonName.FamilyName(),
		)
		So(err, ShouldBeNil)

		Convey("\nSCENARIO 1: Change a Customer's name", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("When ChangeCustomerName", func() {
					recordedEvents := customer.ChangeName(eventStream, changeName)

					Convey("Then CustomerNameChanged", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						nameChanged, ok := recordedEvents[0].(events.CustomerNameChanged)
						So(ok, ShouldBeTrue)
						So(nameChanged, ShouldNotBeNil)
						So(nameChanged.CustomerID().Equals(customerID), ShouldBeTrue)
						So(nameChanged.PersonName().Equals(changedPersonName), ShouldBeTrue)
						So(nameChanged.StreamVersion(), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Try to change a Customer's name to the value he registered with", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("When ChangeCustomerName", func() {
					changeName, err = commands.BuildChangeCustomerName(
						customerID.ID(),
						personName.GivenName(),
						personName.FamilyName(),
					)
					So(err, ShouldBeNil)

					recordedEvents := customer.ChangeName(eventStream, changeName)

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
					nameChanged := events.CustomerNameWasChanged(
						customerID,
						changedPersonName,
						2,
					)

					eventStream = append(eventStream, nameChanged)

					Convey("When ChangeCustomerName", func() {
						recordedEvents := customer.ChangeName(eventStream, changeName)

						Convey("Then no event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})
	})
}
