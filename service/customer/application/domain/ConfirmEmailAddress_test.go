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
	Convey("Given a Customer has an unconfirmed emailAddress", t, func() {
		id, err := values.BuildCustomerID("64bcf656-da30-4f5a-b0b5-aead60965aa3")
		So(err, ShouldBeNil)
		emailAddress, err := values.BuildEmailAddress("john@doe.com")
		So(err, ShouldBeNil)
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName, err := values.BuildPersonName("John", "Doe")
		So(err, ShouldBeNil)

		currentStreamVersion := uint(1)

		eventStream := es.DomainEvents{
			events.ItWasRegistered(id, emailAddress, confirmationHash, personName, currentStreamVersion),
		}

		Convey("When it is confirmed with the right confirmationHash", func() {
			confirmEmailAddress, err := commands.NewConfirmEmailAddress(
				id.ID(),
				emailAddress.EmailAddress(),
				confirmationHash.Hash(),
			)
			So(err, ShouldBeNil)

			recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)

			Convey("It should succeed", func() {
				So(recordedEvents, ShouldHaveLength, 1)
				emailAddressConfirmed, ok := recordedEvents[0].(events.EmailAddressConfirmed)
				So(ok, ShouldBeTrue)
				So(emailAddressConfirmed, ShouldNotBeNil)
				So(emailAddressConfirmed.CustomerID().Equals(id), ShouldBeTrue)
				So(emailAddressConfirmed.EmailAddress().Equals(emailAddress), ShouldBeTrue)
				So(emailAddressConfirmed.StreamVersion(), ShouldEqual, currentStreamVersion+1)

				eventStream = append(eventStream, recordedEvents...)
				currentStreamVersion++

				Convey("And when it is confirmed again", func() {
					recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)

					Convey("It should be ignored", func() {
						So(recordedEvents, ShouldBeEmpty)
					})
				})

				Convey("And when it is changed", func() {
					emailAddress, err := values.BuildEmailAddress("john+different@doe.com")
					So(err, ShouldBeNil)
					confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
					emailAddressChanged := events.EmailAddressWasChanged(id, emailAddress, confirmationHash, currentStreamVersion)
					eventStream = append(eventStream, emailAddressChanged)

					recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)

					Convey("It should be marked as unconfirmed", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						_, ok := recordedEvents[0].(events.EmailAddressConfirmationFailed)
						So(ok, ShouldBeTrue)
					})
				})
			})
		})

		Convey("When the emailAddress is confirmed with some wrong confirmationHash", func() {
			wrongConfirmationHash, err := values.BuildConfirmationHash("some_not_matching_hash")
			So(err, ShouldBeNil)

			confirmEmailAddress, err := commands.NewConfirmEmailAddress(
				id.ID(),
				emailAddress.EmailAddress(),
				wrongConfirmationHash.Hash(),
			)
			So(err, ShouldBeNil)

			recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)

			Convey("It should fail", func() {
				So(recordedEvents, ShouldHaveLength, 1)
				emailAddressConfirmationFailed, ok := recordedEvents[0].(events.EmailAddressConfirmationFailed)
				So(ok, ShouldBeTrue)
				So(emailAddressConfirmationFailed, ShouldNotBeNil)
				So(emailAddressConfirmationFailed.CustomerID().Equals(id), ShouldBeTrue)
				So(emailAddressConfirmationFailed.ConfirmationHash().Equals(wrongConfirmationHash), ShouldBeTrue)
				So(emailAddressConfirmationFailed.StreamVersion(), ShouldEqual, currentStreamVersion+1)
			})
		})
	})
}
