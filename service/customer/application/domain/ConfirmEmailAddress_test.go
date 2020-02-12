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

func TestConfirmEmailAddress(t *testing.T) {
	Convey("Given a Customer with an unconfirmed emailAddress", t, func() {
		id := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("kevin@ball.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("Kevin", "Ball")

		currentStreamVersion := uint(1)
		registered := events.ItWasRegistered(
			id,
			emailAddress,
			confirmationHash,
			personName,
			currentStreamVersion,
		)

		eventStream := es.DomainEvents{registered}

		Convey("When the emailAddress is confirmed with the right confirmationHash", func() {
			confirmEmailAddress, err := commands.NewConfirmEmailAddress(
				registered.CustomerID().ID(),
				registered.EmailAddress().EmailAddress(),
				registered.ConfirmationHash().Hash(),
			)
			So(err, ShouldBeNil)

			recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)

			Convey("Then it should record EmailAddressConfirmed", func() {
				So(recordedEvents, ShouldHaveLength, 1)
				emailAddressConfirmed, ok := recordedEvents[0].(events.EmailAddressConfirmed)
				So(ok, ShouldBeTrue)
				So(emailAddressConfirmed.CustomerID().Equals(registered.CustomerID()), ShouldBeTrue)
				So(emailAddressConfirmed.EmailAddress().Equals(registered.EmailAddress()), ShouldBeTrue)
				So(emailAddressConfirmed.StreamVersion(), ShouldEqual, currentStreamVersion+1)

				eventStream = append(eventStream, emailAddressConfirmed)
				currentStreamVersion++

				Convey("And when it is confirmed again", func() {
					recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)

					Convey("It should be ignored", func() {
						So(recordedEvents, ShouldBeEmpty)
					})
				})
			})
		})

		Convey("When the emailAddress is confirmed with a wrong confirmationHash", func() {
			confirmEmailAddress, err := commands.NewConfirmEmailAddress(
				registered.CustomerID().ID(),
				registered.EmailAddress().EmailAddress(),
				"some_not_matching_hash",
			)
			So(err, ShouldBeNil)

			recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)

			Convey("Then it should record EmailAddressConfirmationFailed", func() {
				So(recordedEvents, ShouldHaveLength, 1)
				emailAddressConfirmationFailed, ok := recordedEvents[0].(events.EmailAddressConfirmationFailed)
				So(ok, ShouldBeTrue)
				So(emailAddressConfirmationFailed.CustomerID().Equals(confirmEmailAddress.CustomerID()), ShouldBeTrue)
				So(emailAddressConfirmationFailed.EmailAddress().Equals(confirmEmailAddress.EmailAddress()), ShouldBeTrue)
				So(emailAddressConfirmationFailed.ConfirmationHash().Equals(confirmEmailAddress.ConfirmationHash()), ShouldBeTrue)
				So(emailAddressConfirmationFailed.StreamVersion(), ShouldEqual, currentStreamVersion+1)
			})
		})
	})
}
