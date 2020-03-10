package customer_test

import (
	"go-iddd/service/customer/application/writemodel/domain/customer"
	"go-iddd/service/customer/application/writemodel/domain/customer/commands"
	"go-iddd/service/customer/application/writemodel/domain/customer/events"
	"go-iddd/service/lib/es"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRegister(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		registerCustomer, err := commands.BuildRegisterCustomer(
			"kevin@ball.com",
			"Kevin",
			"Ball",
		)
		So(err, ShouldBeNil)

		Convey("\nSCENARIO: Register a Customer", func() {
			Convey("When RegisterCustomer", func() {
				recordedEvents := customer.Register(registerCustomer)

				Convey("Then CustomerRegistered", func() {
					ThenCustomerRegistered(recordedEvents, registerCustomer)
				})
			})
		})
	})
}

func ThenCustomerRegistered(recordedEvents es.DomainEvents, register commands.RegisterCustomer) {
	So(recordedEvents, ShouldHaveLength, 1)
	registered, ok := recordedEvents[0].(events.CustomerRegistered)
	So(ok, ShouldBeTrue)
	So(registered.CustomerID().Equals(register.CustomerID()), ShouldBeTrue)
	So(registered.EmailAddress().Equals(register.EmailAddress()), ShouldBeTrue)
	So(registered.ConfirmationHash().Equals(register.ConfirmationHash()), ShouldBeTrue)
	So(registered.PersonName().Equals(register.PersonName()), ShouldBeTrue)
	So(registered.StreamVersion(), ShouldEqual, uint(1))
}
