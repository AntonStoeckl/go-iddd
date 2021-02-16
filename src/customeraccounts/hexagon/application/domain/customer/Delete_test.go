package customer_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDelete(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		customerID := value.GenerateCustomerID()
		emailAddress := value.RebuildEmailAddress("kevin@ball.com")
		confirmationHash := value.GenerateConfirmationHash(emailAddress.String())
		personName := value.RebuildPersonName("Kevin", "Ball")

		command, err := domain.BuildDeleteCustomer(customerID.String())
		So(err, ShouldBeNil)

		customerRegistered := domain.BuildCustomerRegistered(
			customerID,
			emailAddress,
			confirmationHash,
			personName,
			es.GenerateMessageID(),
			1,
		)

		customerDeleted := domain.BuildCustomerDeleted(
			customerID,
			emailAddress,
			es.GenerateMessageID(),
			2,
		)

		Convey("\nSCENARIO 1: Delete a Customer's account", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("When DeleteCustomer", func() {
					recordedEvents := customer.Delete(eventStream, command)

					Convey("Then CustomerDeleted", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						event, ok := recordedEvents[0].(domain.CustomerDeleted)
						So(ok, ShouldBeTrue)
						So(event, ShouldNotBeNil)
						So(event.CustomerID().Equals(customerID), ShouldBeTrue)
						So(event.EmailAddress().Equals(emailAddress), ShouldBeTrue)
						So(event.IsFailureEvent(), ShouldBeFalse)
						So(event.FailureReason(), ShouldBeNil)
						So(event.Meta().CausationID(), ShouldEqual, command.MessageID().String())
						So(event.Meta().StreamVersion(), ShouldEqual, uint(2))
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Try to delete a Customer's account again", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("and CustomerDeleted", func() {
					eventStream = append(eventStream, customerDeleted)

					Convey("When DeleteCustomer", func() {
						recordedEvents := customer.Delete(eventStream, command)

						Convey("Then no Event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})
	})
}
