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
		emailAddress, err := values.BuildEmailAddress("john@doe.com")
		So(err, ShouldBeNil)
		personName, err := values.BuildPersonName("John", "Doe")
		So(err, ShouldBeNil)

		register, err := commands.NewRegister(
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
			So(registered, ShouldNotBeZeroValue)
			So(registered.CustomerID().Equals(register.CustomerID()), ShouldBeTrue)
			So(registered.EmailAddress().Equals(emailAddress), ShouldBeTrue)
			So(registered.ConfirmationHash(), ShouldNotBeZeroValue)
			So(registered.PersonName().Equals(personName), ShouldBeTrue)
			So(registered.StreamVersion(), ShouldEqual, uint(1))
		})
	})
}
