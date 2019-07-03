package domain_test

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/mocks"
	"go-iddd/customer/domain/values"
	"go-iddd/shared"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestChangeEmailAddressOfCustomer(t *testing.T) {
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

		Convey("When an emailAddress is changed", func() {
			newEmailAddress, err := values.NewEmailAddress("john+changed@doe.com")
			So(err, ShouldBeNil)
			changeEmailAddress, err := commands.NewChangeEmailAddress(
				id.String(),
				newEmailAddress.EmailAddress(),
			)
			So(err, ShouldBeNil)

			err = customer.Execute(changeEmailAddress)

			Convey("It should succeed", func() {
				So(err, ShouldBeNil)

				Convey("And it should record that a Customer's emailAddress was changed", func() {
					recordedEvents := customer.RecordedEvents(true)
					emailAddressChanged := mocks.FindCustomerEventIn(
						recordedEvents,
						new(events.EmailAddressChanged),
					).(*events.EmailAddressChanged)

					So(emailAddressChanged, ShouldNotBeNil)
					So(emailAddressChanged.ID().Equals(id), ShouldBeTrue)
					So(emailAddressChanged.EmailAddress().Equals(newEmailAddress), ShouldBeTrue)
					So(emailAddressChanged.StreamVersion(), ShouldEqual, currentStreamVersion+1)

					Convey("And it should not record anything else", func() {
						So(recordedEvents, ShouldHaveLength, 1)
					})

					Convey("And when it is changed to the same value again", func() {
						err = customer.Execute(changeEmailAddress)

						Convey("It should be ignored", func() {
							So(err, ShouldBeNil)
							recordedEvents := customer.RecordedEvents(false)
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})
	})
}
