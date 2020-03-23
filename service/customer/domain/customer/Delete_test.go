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
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("Kevin", "Ball")

		customerWasRegistered := events.CustomerWasRegistered(
			customerID,
			emailAddress,
			confirmationHash,
			personName,
			1,
		)

		Convey("\nSCENARIO 1: Delete a Customer's account", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("When DeleteCustomer", func() {
					recordedEvents := customer.Delete(eventStream)

					Convey("Then CustomerDeleted", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						nameChanged, ok := recordedEvents[0].(events.CustomerDeleted)
						So(ok, ShouldBeTrue)
						So(nameChanged, ShouldNotBeNil)
						So(nameChanged.CustomerID().Equals(customerID), ShouldBeTrue)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Try to delete a Customer's account again", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("and CustomerDeleted", func() {
					customerDeleted := events.CustomerWasDeleted(customerID, 2)
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
