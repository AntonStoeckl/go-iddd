package customer_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/service/shared/es"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDelete(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		customerID := value.GenerateCustomerID()
		emailAddress := value.RebuildEmailAddress("kevin@ball.com")
		confirmationHash := value.GenerateConfirmationHash(emailAddress.String())
		personName := value.RebuildPersonName("Kevin", "Ball")

		customerWasRegistered := domain.BuildCustomerRegistered(
			customerID,
			emailAddress,
			confirmationHash,
			personName,
			1,
		)

		deleteCmd := domain.BuildDeleteCustomer(customerID)

		Convey("\nSCENARIO 1: Delete a Customer's account", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerWasRegistered}

				Convey("When DeleteCustomer", func() {
					recordedEvents := customer.Delete(eventStream, deleteCmd)

					Convey("Then CustomerDeleted", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						customerDeleted, ok := recordedEvents[0].(domain.CustomerDeleted)
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
					customerDeleted := domain.BuildCustomerDeleted(customerID, emailAddress, 2)
					eventStream = append(eventStream, customerDeleted)

					Convey("When DeleteCustomer", func() {
						recordedEvents := customer.Delete(eventStream, deleteCmd)

						Convey("Then no Event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})
	})
}
