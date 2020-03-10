package domain_test

import (
	"go-iddd/service/customer/application/domain"
	"go-iddd/service/customer/application/domain/commands"
	"go-iddd/service/customer/application/domain/events"
	"go-iddd/service/customer/application/domain/values"
	"go-iddd/service/lib/es"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestConfirmCustomerEmailAddress(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
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

		confirmEmailAddress, err := commands.NewConfirmEmailAddress(
			customerID.ID(),
			confirmationHash.Hash(),
		)
		So(err, ShouldBeNil)

		confirmEmailAddressWithInvalidHash, err := commands.NewConfirmEmailAddress(
			customerID.ID(),
			invalidConfirmationHash.Hash(),
		)
		So(err, ShouldBeNil)

		Convey("\nSCENARIO 1: Confirm a Customer's emailAddress with the right confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("When ConfirmCustomerEmailAddress", func() {
					recordedEvents := domain.ConfirmCustomerEmailAddress(eventStream, confirmEmailAddress)

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
					recordedEvents := domain.ConfirmCustomerEmailAddress(eventStream, confirmEmailAddressWithInvalidHash)

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
						recordedEvents := domain.ConfirmCustomerEmailAddress(eventStream, confirmEmailAddress)

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
						recordedEvents := domain.ConfirmCustomerEmailAddress(eventStream, confirmEmailAddressWithInvalidHash)

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
	})
}
