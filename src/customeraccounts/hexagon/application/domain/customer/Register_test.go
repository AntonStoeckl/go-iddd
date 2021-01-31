package customer_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRegister(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		var err error

		emailAddress, err := value.BuildEmailAddress("kevin@ball.com")
		So(err, ShouldBeNil)

		personName, err := value.BuildPersonName("Kevin", "Ball")
		So(err, ShouldBeNil)

		register := domain.BuildRegisterCustomer(
			value.GenerateCustomerID(),
			emailAddress,
			value.GenerateConfirmationHash("kevin@ball.com"),
			personName,
		)

		Convey("\nSCENARIO: Register a Customer", func() {
			Convey("When RegisterCustomer", func() {
				registered := customer.Register(register)

				Convey("Then CustomerRegistered", func() {
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
