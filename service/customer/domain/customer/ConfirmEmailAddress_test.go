package customer_test

import (
	"testing"

	"github.com/AntonStoeckl/go-iddd/service/lib"
	"github.com/cockroachdb/errors"

	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/commands"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/events"
	"github.com/AntonStoeckl/go-iddd/service/customer/domain/customer/values"
	"github.com/AntonStoeckl/go-iddd/service/lib/es"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConfirmEmailAddress(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		var err error
		var recordedEvents es.DomainEvents

		customerID := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("kevin@ball.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		invalidConfirmationHash := values.RebuildConfirmationHash("invalid_hash")
		personName := values.RebuildPersonName("Kevin", "Ball")

		customerWasRegistered := events.CustomerWasRegistered(
			customerID,
			emailAddress,
			confirmationHash,
			personName,
			1,
		)

		customerEmailAddressWasConfirmed := events.CustomerEmailAddressWasConfirmed(
			customerID,
			emailAddress,
			2,
		)

		confirmEmailAddress, err := commands.BuildConfirmCustomerEmailAddress(
			customerID.ID(),
			confirmationHash.Hash(),
		)
		So(err, ShouldBeNil)

		confirmEmailAddressWithInvalidHash, err := commands.BuildConfirmCustomerEmailAddress(
			customerID.ID(),
			invalidConfirmationHash.Hash(),
		)
		So(err, ShouldBeNil)

		Convey("\nSCENARIO 1: Confirm a Customer's emailAddress with the right confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("When ConfirmCustomerEmailAddress", func() {
					recordedEvents, err = customer.ConfirmEmailAddress(eventStream, confirmEmailAddress)
					So(err, ShouldBeNil)

					Convey("Then CustomerEmailAddressConfirmed", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						emailAddressConfirmed, ok := recordedEvents[0].(events.CustomerEmailAddressConfirmed)
						So(ok, ShouldBeTrue)
						So(emailAddressConfirmed.CustomerID().Equals(customerID), ShouldBeTrue)
						So(emailAddressConfirmed.EmailAddress().Equals(emailAddress), ShouldBeTrue)
						So(emailAddressConfirmed.StreamVersion(), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Confirm a Customer's emailAddress with a wrong confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("When ConfirmCustomerEmailAddress", func() {
					recordedEvents, err = customer.ConfirmEmailAddress(eventStream, confirmEmailAddressWithInvalidHash)
					So(err, ShouldBeNil)

					Convey("Then CustomerEmailAddressConfirmationFailed", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						emailAddressConfirmationFailed, ok := recordedEvents[0].(events.CustomerEmailAddressConfirmationFailed)
						So(ok, ShouldBeTrue)
						So(emailAddressConfirmationFailed.CustomerID().Equals(customerID), ShouldBeTrue)
						So(emailAddressConfirmationFailed.EmailAddress().Equals(emailAddress), ShouldBeTrue)
						So(emailAddressConfirmationFailed.ConfirmationHash().Equals(invalidConfirmationHash), ShouldBeTrue)
						So(emailAddressConfirmationFailed.StreamVersion(), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Try to confirm a Customer's emailAddress again with the right confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

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
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("and CustomerEmailAddressConfirmed", func() {
					eventStream = append(eventStream, customerEmailAddressWasConfirmed)

					Convey("When ConfirmCustomerEmailAddress", func() {
						recordedEvents, err = customer.ConfirmEmailAddress(eventStream, confirmEmailAddressWithInvalidHash)
						So(err, ShouldBeNil)

						Convey("Then CustomerEmailAddressConfirmationFailed", func() {
							So(recordedEvents, ShouldHaveLength, 1)
							emailAddressConfirmationFailed, ok := recordedEvents[0].(events.CustomerEmailAddressConfirmationFailed)
							So(ok, ShouldBeTrue)
							So(emailAddressConfirmationFailed.CustomerID().Equals(customerID), ShouldBeTrue)
							So(emailAddressConfirmationFailed.EmailAddress().Equals(emailAddress), ShouldBeTrue)
							So(emailAddressConfirmationFailed.ConfirmationHash().Equals(invalidConfirmationHash), ShouldBeTrue)
							So(emailAddressConfirmationFailed.StreamVersion(), ShouldEqual, 3)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 5: Try to confirm a Customer's emailAddress when the account was deleted", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("Given CustomerDeleted", func() {
					eventStream = append(
						eventStream,
						events.CustomerWasDeleted(customerID, 2),
					)

					Convey("When ConfirmCustomerEmailAddress", func() {
						_, err := customer.ConfirmEmailAddress(eventStream, confirmEmailAddress)

						Convey("Then it should report an error", func() {
							So(err, ShouldBeError)
							So(errors.Is(err, lib.ErrNotFound), ShouldBeTrue)
						})
					})
				})
			})
		})
	})
}
