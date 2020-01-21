package domain_test

import (
	"go-iddd/service/customer/application/domain"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/application/domain/values"
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
			So(registered.ConfirmationHash().Equals(register.ConfirmationHash()), ShouldBeTrue)
			So(registered.PersonName().Equals(personName), ShouldBeTrue)
			So(registered.StreamVersion(), ShouldEqual, uint(1))
		})
	})
}
