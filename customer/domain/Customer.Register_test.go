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
		id, err := values.BuildCustomerID("64bcf656-da30-4f5a-b0b5-aead60965aa3")
		So(err, ShouldBeNil)
		emailAddress, err := values.BuildEmailAddress("john@doe.com")
		So(err, ShouldBeNil)
		personName, err := values.BuildPersonName("John", "Doe")
		So(err, ShouldBeNil)

		register, err := commands.NewRegister(
			id.ID(),
			emailAddress.EmailAddress(),
			personName.GivenName(),
			personName.FamilyName(),
		)
		So(err, ShouldBeNil)

		recordedEvents := domain.RegisterCustomer(register)

		Convey("It should record CustomerRegistered", func() {
			So(recordedEvents, ShouldHaveLength, 1)
			registered, ok := recordedEvents[0].(events.Registered)
			So(ok, ShouldBeTrue)
			So(registered, ShouldNotBeNil)
			So(registered.CustomerID().Equals(id), ShouldBeTrue)
			So(registered.EmailAddress().Equals(emailAddress), ShouldBeTrue)
			_, err := values.BuildConfirmationHash(registered.ConfirmationHash().Hash())
			So(err, ShouldBeNil)
			So(registered.PersonName().Equals(personName), ShouldBeTrue)
			So(registered.StreamVersion(), ShouldEqual, uint(1))
		})
	})
}
