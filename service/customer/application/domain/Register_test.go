package domain_test

import (
	"go-iddd/service/customer/application/domain"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegister(t *testing.T) {
	Convey("When a Customer is registered", t, func() {
		register, err := commands.NewRegister(
			"fiona@galagher.com",
			"Fiona",
			"Galagher",
		)
		So(err, ShouldBeNil)

		recordedEvents := domain.RegisterCustomer(register)

		Convey("It should record CustomerRegistered", func() {
			So(recordedEvents, ShouldHaveLength, 1)
			registered, ok := recordedEvents[0].(events.Registered)
			So(ok, ShouldBeTrue)
			So(registered.CustomerID().Equals(register.CustomerID()), ShouldBeTrue)
			So(registered.EmailAddress().Equals(register.EmailAddress()), ShouldBeTrue)
			So(registered.ConfirmationHash().Equals(register.ConfirmationHash()), ShouldBeTrue)
			So(registered.PersonName().Equals(register.PersonName()), ShouldBeTrue)
			So(registered.StreamVersion(), ShouldEqual, uint(1))
		})
	})
}
