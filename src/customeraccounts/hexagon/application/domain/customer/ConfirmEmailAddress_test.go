package customer_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer"
	"github.com/AntonStoeckl/go-iddd/src/customeraccounts/hexagon/application/domain/customer/value"
	"github.com/AntonStoeckl/go-iddd/src/shared"
	"github.com/AntonStoeckl/go-iddd/src/shared/es"
	"github.com/cockroachdb/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConfirmEmailAddress(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		var err error
		var recordedEvents es.RecordedEvents

		customerID := value.GenerateCustomerID()
		emailAddress := value.RebuildEmailAddress("kevin@ball.com")
		confirmationHash := value.GenerateConfirmationHash(emailAddress.String())
		invalidConfirmationHash := value.RebuildConfirmationHash("invalid_hash")
		personName := value.RebuildPersonName("Kevin", "Ball")

		customerWasRegistered := domain.BuildCustomerRegistered(
			customerID,
			emailAddress,
			confirmationHash,
			personName,
			1,
		)

		customerEmailAddressWasConfirmed := domain.BuildCustomerEmailAddressConfirmed(
			customerID,
			emailAddress,
			2,
		)

		confirmEmailAddress := domain.BuildConfirmCustomerEmailAddress(
			customerID,
			confirmationHash,
		)

		confirmEmailAddressWithInvalidHash := domain.BuildConfirmCustomerEmailAddress(
			customerID,
			invalidConfirmationHash,
		)

		Convey("\nSCENARIO 1: Confirm a Customer's emailAddress with the right confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerWasRegistered}

				Convey("When ConfirmCustomerEmailAddress", func() {
					recordedEvents, err = customer.ConfirmEmailAddress(eventStream, confirmEmailAddress)
					So(err, ShouldBeNil)

					Convey("Then CustomerEmailAddressConfirmed", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						emailAddressConfirmed, ok := recordedEvents[0].(domain.CustomerEmailAddressConfirmed)
						So(ok, ShouldBeTrue)
						So(emailAddressConfirmed.CustomerID().Equals(customerID), ShouldBeTrue)
						So(emailAddressConfirmed.EmailAddress().Equals(emailAddress), ShouldBeTrue)
						So(emailAddressConfirmed.Meta().StreamVersion(), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Confirm a Customer's emailAddress with a wrong confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerWasRegistered}

				Convey("When ConfirmCustomerEmailAddress", func() {
					recordedEvents, err = customer.ConfirmEmailAddress(eventStream, confirmEmailAddressWithInvalidHash)
					So(err, ShouldBeNil)

					Convey("Then CustomerEmailAddressConfirmationFailed", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						emailAddressConfirmationFailed, ok := recordedEvents[0].(domain.CustomerEmailAddressConfirmationFailed)
						So(ok, ShouldBeTrue)
						So(emailAddressConfirmationFailed.CustomerID().Equals(customerID), ShouldBeTrue)
						So(emailAddressConfirmationFailed.EmailAddress().Equals(emailAddress), ShouldBeTrue)
						So(emailAddressConfirmationFailed.ConfirmationHash().Equals(invalidConfirmationHash), ShouldBeTrue)
						So(emailAddressConfirmationFailed.IsFailureEvent(), ShouldBeTrue)
						So(emailAddressConfirmationFailed.FailureReason(), ShouldBeError)
						So(emailAddressConfirmationFailed.Meta().StreamVersion(), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Try to confirm a Customer's emailAddress again with the right confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerWasRegistered}

				Convey("and CustomerEmailAddressConfirmed", func() {
					eventStream = append(eventStream, customerEmailAddressWasConfirmed)

					Convey("When ConfirmCustomerEmailAddress", func() {
						recordedEvents, err = customer.ConfirmEmailAddress(eventStream, confirmEmailAddress)
						So(err, ShouldBeNil)

						Convey("Then no event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 4: Try to confirm a Customer's emailAddress again with a wrong confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerWasRegistered}

				Convey("and CustomerEmailAddressConfirmed", func() {
					eventStream = append(eventStream, customerEmailAddressWasConfirmed)

					Convey("When ConfirmCustomerEmailAddress", func() {
						recordedEvents, err = customer.ConfirmEmailAddress(eventStream, confirmEmailAddressWithInvalidHash)
						So(err, ShouldBeNil)

						Convey("Then CustomerEmailAddressConfirmationFailed", func() {
							So(recordedEvents, ShouldHaveLength, 1)
							emailAddressConfirmationFailed, ok := recordedEvents[0].(domain.CustomerEmailAddressConfirmationFailed)
							So(ok, ShouldBeTrue)
							So(emailAddressConfirmationFailed.CustomerID().Equals(customerID), ShouldBeTrue)
							So(emailAddressConfirmationFailed.EmailAddress().Equals(emailAddress), ShouldBeTrue)
							So(emailAddressConfirmationFailed.ConfirmationHash().Equals(invalidConfirmationHash), ShouldBeTrue)
							So(emailAddressConfirmationFailed.IsFailureEvent(), ShouldBeTrue)
							So(emailAddressConfirmationFailed.FailureReason(), ShouldBeError)
							So(emailAddressConfirmationFailed.Meta().StreamVersion(), ShouldEqual, 3)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 5: Try to confirm a Customer's emailAddress when the account was deleted", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerWasRegistered}

				Convey("Given CustomerDeleted", func() {
					eventStream = append(
						eventStream,
						domain.BuildCustomerDeleted(customerID, emailAddress, 2),
					)

					Convey("When ConfirmCustomerEmailAddress", func() {
						_, err := customer.ConfirmEmailAddress(eventStream, confirmEmailAddress)

						Convey("Then it should report an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, shared.ErrNotFound), ShouldBeTrue)
						})
					})
				})
			})
		})
	})
}
