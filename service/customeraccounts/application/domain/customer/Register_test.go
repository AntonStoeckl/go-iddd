package customer_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/application/domain/customer"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRegister(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		register, err := domain.BuildRegisterCustomer(
			"kevin@ball.com",
			"Kevin",
			"Ball",
		)
		So(err, ShouldBeNil)

		Convey("\nSCENARIO: Register a Customer", func() {
			Convey("When RegisterCustomer", func() {
				recordedEvents := customer.Register(register)

				Convey("Then CustomerRegistered", func() {
					So(recordedEvents, ShouldHaveLength, 1)
					registered, ok := recordedEvents[0].(domain.CustomerRegistered)
					So(ok, ShouldBeTrue)
					So(registered.CustomerID().Equals(register.CustomerID()), ShouldBeTrue)
					So(registered.EmailAddress().Equals(register.EmailAddress()), ShouldBeTrue)
					So(registered.ConfirmationHash().Equals(register.ConfirmationHash()), ShouldBeTrue)
					So(registered.PersonName().Equals(register.PersonName()), ShouldBeTrue)
					So(registered.IsFailureEvent(), ShouldBeFalse)
					So(registered.FailureReason(), ShouldBeNil)
					So(registered.Meta().StreamVersion(), ShouldEqual, uint(1))
				})
			})
		})
	})
}
