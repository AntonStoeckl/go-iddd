package domain_test

import (
	"go-iddd/service/customer/application/domain"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib/es"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChangeEmailAddress(t *testing.T) {
	Convey("Given a Customer with a confirmed emailAddress", t, func() {
		id := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("john@doe.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("John", "Doe")

		currentStreamVersion := uint(1)
		registered := events.ItWasRegistered(id, emailAddress, confirmationHash, personName, currentStreamVersion)

		currentStreamVersion++
		emailAddressConfirmed := events.EmailAddressWasConfirmed(
			registered.CustomerID(),
			registered.EmailAddress(),
			currentStreamVersion,
		)

		eventStream := es.DomainEvents{registered, emailAddressConfirmed}

		Convey("When the emailAddress is changed", func() {
			changeEmailAddress, err := commands.NewChangeEmailAddress(
				registered.CustomerID().ID(),
				"john+changed@doe.com",
			)
			So(err, ShouldBeNil)

			recordedEvents := domain.ChangeEmailAddress(eventStream, changeEmailAddress)

			Convey("Then it should record EmailAddressChanged", func() {
				So(recordedEvents, ShouldHaveLength, 1)
				emailAddressChanged, ok := recordedEvents[0].(events.EmailAddressChanged)
				So(ok, ShouldBeTrue)
				So(emailAddressChanged, ShouldNotBeNil)
				So(emailAddressChanged.CustomerID().Equals(changeEmailAddress.CustomerID()), ShouldBeTrue)
				So(emailAddressChanged.EmailAddress().Equals(changeEmailAddress.EmailAddress()), ShouldBeTrue)
				So(emailAddressChanged.ConfirmationHash().Equals(changeEmailAddress.ConfirmationHash()), ShouldBeTrue)
				So(emailAddressChanged.StreamVersion(), ShouldEqual, currentStreamVersion+1)

				eventStream = append(eventStream, emailAddressChanged)
				currentStreamVersion++

				Convey("And when it is changed to the same value again", func() {
					recordedEvents := domain.ChangeEmailAddress(eventStream, changeEmailAddress)

					Convey("Then it should be ignored", func() {
						So(recordedEvents, ShouldBeEmpty)
					})
				})

				Convey("And when the changed emailAddress is confirmed with the right confirmationHash", func() {
					confirmEmailAddress, err := commands.NewConfirmEmailAddress(
						emailAddressChanged.CustomerID().ID(),
						emailAddressChanged.EmailAddress().EmailAddress(),
						emailAddressChanged.ConfirmationHash().Hash(),
					)
					So(err, ShouldBeNil)

					recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)

					Convey("Then it should record EmailAddressConfirmed", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						emailAddressConfirmed, ok := recordedEvents[0].(events.EmailAddressConfirmed)
						So(ok, ShouldBeTrue)
						So(emailAddressConfirmed.CustomerID().Equals(confirmEmailAddress.CustomerID()), ShouldBeTrue)
						So(emailAddressConfirmed.EmailAddress().Equals(confirmEmailAddress.EmailAddress()), ShouldBeTrue)
						So(emailAddressConfirmed.StreamVersion(), ShouldEqual, currentStreamVersion+1)
					})
				})
			})
		})
	})
}
