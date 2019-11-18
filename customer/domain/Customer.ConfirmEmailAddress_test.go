package domain_test

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfirmEmailAddressOfCustomer(t *testing.T) {
	Convey("Given a Customer", t, func() {
		id, err := values.CustomerIDFrom("64bcf656-da30-4f5a-b0b5-aead60965aa3")
		So(err, ShouldBeNil)
		emailAddress, err := values.EmailAddressFrom("john@doe.com")
		So(err, ShouldBeNil)
		confirmableEmailAddress := emailAddress.ToConfirmable()
		personName, err := values.PersonNameFrom("John", "Doe")
		So(err, ShouldBeNil)

		currentStreamVersion := uint(1)

		customer, err := domain.ReconstituteCustomerFrom(
			shared.DomainEvents{
				events.ItWasRegistered(id, confirmableEmailAddress, personName, currentStreamVersion),
			},
		)
		So(err, ShouldBeNil)

		Convey("When an unconfirmed emailAddress is confirmed", func() {
			confirmEmailAddress, err := commands.NewConfirmEmailAddress(
				id.String(),
				emailAddress.EmailAddress(),
				confirmableEmailAddress.ConfirmationHash(),
			)
			So(err, ShouldBeNil)

			recordedEvents := customer.ConfirmEmailAddress(confirmEmailAddress)

			Convey("And it should record EmailAddressConfirmed", func() {
				So(recordedEvents, ShouldHaveLength, 1)
				emailAddressConfirmed, ok := recordedEvents[0].(*events.EmailAddressConfirmed)
				So(ok, ShouldBeTrue)
				So(emailAddressConfirmed, ShouldNotBeNil)
				So(emailAddressConfirmed.CustomerID().Equals(id), ShouldBeTrue)
				So(emailAddressConfirmed.EmailAddress().Equals(emailAddress), ShouldBeTrue)
				So(emailAddressConfirmed.StreamVersion(), ShouldEqual, currentStreamVersion+1)

				Convey("And when it is confirmed again", func() {
					recordedEvents := customer.ConfirmEmailAddress(confirmEmailAddress)

					Convey("It should be ignored", func() {
						So(recordedEvents, ShouldBeEmpty)
					})
				})
			})
		})

		Convey("When an emailAddress is confirmed with a wrong confirmationHash", func() {
			wrongConfirmationHash, err := values.ConfirmationHashFrom("some_not_matching_hash")
			So(err, ShouldBeNil)

			confirmEmailAddress, err := commands.NewConfirmEmailAddress(
				id.String(),
				emailAddress.EmailAddress(),
				wrongConfirmationHash.Hash(),
			)
			So(err, ShouldBeNil)

			recordedEvents := customer.ConfirmEmailAddress(confirmEmailAddress)

			Convey("It should record EmailAddressConfirmationFailed", func() {
				So(recordedEvents, ShouldHaveLength, 1)
				emailAddressConfirmationFailed, ok := recordedEvents[0].(*events.EmailAddressConfirmationFailed)
				So(ok, ShouldBeTrue)
				So(emailAddressConfirmationFailed, ShouldNotBeNil)
				So(emailAddressConfirmationFailed.CustomerID().Equals(id), ShouldBeTrue)
				So(emailAddressConfirmationFailed.ConfirmationHash().Equals(wrongConfirmationHash), ShouldBeTrue)
				So(emailAddressConfirmationFailed.StreamVersion(), ShouldEqual, currentStreamVersion+1)
			})
		})
	})
}
