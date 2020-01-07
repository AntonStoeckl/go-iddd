package domain_test

import (
	"go-iddd/customer/domain"
	"go-iddd/customer/domain/commands"
	"go-iddd/customer/domain/events"
	"go-iddd/customer/domain/values"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegisterCustomer(t *testing.T) {
	Convey("When a Customer is registered", t, func() {
		id, err := values.CustomerIDFrom("64bcf656-da30-4f5a-b0b5-aead60965aa3")
		So(err, ShouldBeNil)
		emailAddress, err := values.EmailAddressFrom("john@doe.com")
		So(err, ShouldBeNil)
		personName, err := values.PersonNameFrom("John", "Doe")
		So(err, ShouldBeNil)

		register, err := commands.NewRegister(
			id.String(),
			emailAddress.EmailAddress(),
			personName.GivenName(),
			personName.FamilyName(),
		)
		So(err, ShouldBeNil)

		recordedEvents := domain.RegisterCustomer(register)

		Convey("It should record CustomerRegistered", func() {
			So(recordedEvents, ShouldHaveLength, 1)
			registered, ok := recordedEvents[0].(*events.Registered)
			So(ok, ShouldBeTrue)
			So(registered, ShouldNotBeNil)
			So(registered.CustomerID().Equals(id), ShouldBeTrue)
			So(registered.ConfirmableEmailAddress().Equals(emailAddress), ShouldBeTrue)
			_, err := values.ConfirmationHashFrom(registered.ConfirmableEmailAddress().ConfirmationHash())
			So(err, ShouldBeNil)
			So(registered.PersonName().Equals(personName), ShouldBeTrue)
			So(registered.StreamVersion(), ShouldEqual, uint(1))
		})
	})
}
