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
		customerID := value.GenerateCustomerID()
		emailAddress := value.RebuildEmailAddress("kevin@ball.com")
		personName := value.RebuildPersonName("Kevin", "Ball")

		command := domain.BuildRegisterCustomer(
			customerID,
			emailAddress,
			personName,
		)

		Convey("\nSCENARIO: Register a Customer", func() {
			Convey("When RegisterCustomer", func() {
				event := customer.Register(command)

				Convey("Then CustomerRegistered", func() {
					So(event.CustomerID().Equals(command.CustomerID()), ShouldBeTrue)
					So(event.EmailAddress().Equals(command.EmailAddress()), ShouldBeTrue)
					So(event.ConfirmationHash().Equals(command.ConfirmationHash()), ShouldBeTrue)
					So(event.PersonName().Equals(command.PersonName()), ShouldBeTrue)
					So(event.IsFailureEvent(), ShouldBeFalse)
					So(event.FailureReason(), ShouldBeNil)
					So(event.Meta().CausationID(), ShouldEqual, command.MessageID().String())
					So(event.Meta().StreamVersion(), ShouldEqual, uint(1))
				})
			})
		})
	})
}
