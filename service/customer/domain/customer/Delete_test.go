package customer_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDelete(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		customerID := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("kevin@ball.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.String())
		personName := values.RebuildPersonName("Kevin", "Ball")

		customerWasRegistered := events.BuildCustomerRegistered(
			customerID,
			emailAddress,
			confirmationHash,
			personName,
			1,
		)

		Convey("\nSCENARIO 1: PurgeCustomerEventStream a Customer's account", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerWasRegistered}

				Convey("When DeleteCustomer", func() {
					recordedEvents := customer.Delete(eventStream)

					Convey("Then CustomerDeleted", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						customerDeleted, ok := recordedEvents[0].(events.CustomerDeleted)
						So(ok, ShouldBeTrue)
						So(customerDeleted, ShouldNotBeNil)
						So(customerDeleted.CustomerID().Equals(customerID), ShouldBeTrue)
						So(customerDeleted.EmailAddress().Equals(emailAddress), ShouldBeTrue)
						So(customerDeleted.IsFailureEvent(), ShouldBeFalse)
						So(customerDeleted.FailureReason(), ShouldBeNil)
						So(customerDeleted.Meta().StreamVersion(), ShouldEqual, uint(2))
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Try to delete a Customer's account again", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerWasRegistered}

				Convey("and CustomerDeleted", func() {
					customerDeleted := events.BuildCustomerDeleted(customerID, emailAddress, 2)
					eventStream = append(eventStream, customerDeleted)

					Convey("When DeleteCustomer", func() {
						recordedEvents := customer.Delete(eventStream)

						Convey("Then no Event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})
	})
}
