package domain_test

import (
	"go-iddd/service/customer/domain"
	"go-iddd/service/customer/domain/commands"
	"go-iddd/service/customer/domain/events"
	"go-iddd/service/customer/domain/values"
	"go-iddd/service/lib"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChangeEmailAddressOfCustomer(t *testing.T) {
	Convey("Given a Customer", t, func() {
		id, err := values.BuildCustomerID("64bcf656-da30-4f5a-b0b5-aead60965aa3")
		So(err, ShouldBeNil)
		emailAddress, err := values.BuildEmailAddress("john@doe.com")
		So(err, ShouldBeNil)
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName, err := values.BuildPersonName("John", "Doe")
		So(err, ShouldBeNil)

		currentStreamVersion := uint(1)

		eventStream := lib.DomainEvents{
			events.ItWasRegistered(id, emailAddress, confirmationHash, personName, currentStreamVersion),
		}

		Convey("When an emailAddress is changed", func() {
			newEmailAddress, err := values.BuildEmailAddress("john+changed@doe.com")
			So(err, ShouldBeNil)

			changeEmailAddress, err := commands.NewChangeEmailAddress(
				id.ID(),
				newEmailAddress.EmailAddress(),
			)
			So(err, ShouldBeNil)

			recordedEvents := domain.ChangeEmailAddress(eventStream, changeEmailAddress)

			Convey("It should record EmailAddressChanged", func() {
				So(recordedEvents, ShouldHaveLength, 1)
				emailAddressChanged, ok := recordedEvents[0].(events.EmailAddressChanged)
				So(ok, ShouldBeTrue)
				So(emailAddressChanged, ShouldNotBeNil)
				So(emailAddressChanged.CustomerID().Equals(id), ShouldBeTrue)
				So(emailAddressChanged.EmailAddress().Equals(newEmailAddress), ShouldBeTrue)
				_, err := values.BuildConfirmationHash(emailAddressChanged.ConfirmationHash().Hash())
				So(err, ShouldBeNil)
				So(emailAddressChanged.StreamVersion(), ShouldEqual, currentStreamVersion+1)

				Convey("And when it is changed to the same value again", func() {
					eventStream = append(eventStream, recordedEvents...)
					recordedEvents := domain.ChangeEmailAddress(eventStream, changeEmailAddress)

					Convey("It should be ignored", func() {
						So(recordedEvents, ShouldBeEmpty)
					})
				})
			})
		})
	})
}
