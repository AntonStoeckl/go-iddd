package domain_test

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
)

func TestConfirmEmailAddressOfCustomer(t *testing.T) {
	Convey("Given a Customer", t, func() {
		id, err := values.RebuildID("64bcf656-da30-4f5a-b0b5-aead60965aa3")
		So(err, ShouldBeNil)

		emailAddress, err := values.NewEmailAddress("john@doe.com")
		So(err, ShouldBeNil)

		confirmableEmailAddress := emailAddress.ToConfirmable()

		personName, err := values.NewPersonName("John", "Doe")
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

			err = customer.Apply(confirmEmailAddress)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)

				Convey("And it should record that a Customer was registered", func() {
					recordedEvents := customer.RecordedEvents()
					emailAddressConfirmed := findCustomerEventIn(
						recordedEvents,
						new(events.EmailAddressConfirmed),
					).(*events.EmailAddressConfirmed)

					So(emailAddressConfirmed, ShouldNotBeNil)
					So(emailAddressConfirmed.ID().Equals(id), ShouldBeTrue)
					So(emailAddressConfirmed.EmailAddress().Equals(emailAddress), ShouldBeTrue)
					So(emailAddressConfirmed.StreamVersion(), ShouldEqual, currentStreamVersion+1)

					Convey("And it should not record anything else", func() {
						So(recordedEvents, ShouldHaveLength, 1)
					})

					Convey("And when it is confirmed again", func() {
						err = customer.Apply(confirmEmailAddress)

						Convey("It should be ignored", func() {
							So(err, ShouldBeNil)
							recordedEvents := customer.RecordedEvents()
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})

		Convey("When a confirmableEmailAddress is confirmed with a wrong emailAddress", func() {
			confirmEmailAddress, err := commands.NewConfirmEmailAddress(
				id.String(),
				"john+outdated@doe.com",
				confirmableEmailAddress.ConfirmationHash(),
			)
			So(err, ShouldBeNil)

			err = customer.Apply(confirmEmailAddress)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrDomainConstraintsViolation), ShouldBeTrue)
			})
		})

		Convey("When a confirmableEmailAddress is confirmed with a wrong confirmationHash", func() {
			confirmEmailAddress, err := commands.NewConfirmEmailAddress(
				id.String(),
				emailAddress.EmailAddress(),
				"some_not_matching_hash",
			)
			So(err, ShouldBeNil)

			err = customer.Apply(confirmEmailAddress)

			Convey("It should fail", func() {
				So(err, ShouldBeError)
				So(xerrors.Is(err, shared.ErrDomainConstraintsViolation), ShouldBeTrue)
			})
		})
	})
}
