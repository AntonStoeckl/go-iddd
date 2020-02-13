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

func TestConfirmEmailAddress(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		id := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("kevin@ball.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("Kevin", "Ball")

		registered := events.ItWasRegistered(
			id,
			emailAddress,
			confirmationHash,
			personName,
			1,
		)

		emailAddressConfirmed := events.EmailAddressWasConfirmed(
			registered.CustomerID(),
			registered.EmailAddress(),
			2,
		)

		confirmEmailAddress, err := commands.NewConfirmEmailAddress(
			registered.CustomerID().ID(),
			registered.EmailAddress().EmailAddress(),
			registered.ConfirmationHash().Hash(),
		)
		So(err, ShouldBeNil)

		confirmEmailAddressWithInvalidHash, err := commands.NewConfirmEmailAddress(
			registered.CustomerID().ID(),
			registered.EmailAddress().EmailAddress(),
			"invalid_hash",
		)
		So(err, ShouldBeNil)

		Convey("\nSCENARIO 1: Confirm a Customer's emailAddress with the right confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{registered}

				Convey("When ConfirmEmailAddress", func() {
					recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)

					Convey("Then EmailAddressConfirmed", func() {
						ThenEmailAddressConfirmed(recordedEvents, confirmEmailAddress)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Confirm a Customer's emailAddress with a wrong confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{registered}

				Convey("When ConfirmEmailAddress", func() {
					recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddressWithInvalidHash)

					Convey("Then EmailAddressConfirmationFailed", func() {
						ThenEmailAddressConfirmationFailed(recordedEvents, confirmEmailAddressWithInvalidHash)
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Try to Confirm a Customer's emailAddress twice with the right confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{registered}

				Convey("and EmailAddressConfirmed", func() {
					eventStream = append(eventStream, emailAddressConfirmed)

					Convey("When ConfirmEmailAddress", func() {
						recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)

						Convey("Then no event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})
	})
}

func ThenEmailAddressConfirmed(recordedEvents es.DomainEvents, confirmEmailAddress commands.ConfirmEmailAddress) {
	So(recordedEvents, ShouldHaveLength, 1)
	emailAddressConfirmed, ok := recordedEvents[0].(events.EmailAddressConfirmed)
	So(ok, ShouldBeTrue)
	So(emailAddressConfirmed.CustomerID().Equals(confirmEmailAddress.CustomerID()), ShouldBeTrue)
	So(emailAddressConfirmed.EmailAddress().Equals(confirmEmailAddress.EmailAddress()), ShouldBeTrue)
	So(emailAddressConfirmed.StreamVersion(), ShouldEqual, 2)
}

func ThenEmailAddressConfirmationFailed(recordedEvents es.DomainEvents, confirmEmailAddress commands.ConfirmEmailAddress) {
	So(recordedEvents, ShouldHaveLength, 1)
	emailAddressConfirmationFailed, ok := recordedEvents[0].(events.EmailAddressConfirmationFailed)
	So(ok, ShouldBeTrue)
	So(emailAddressConfirmationFailed.CustomerID().Equals(confirmEmailAddress.CustomerID()), ShouldBeTrue)
	So(emailAddressConfirmationFailed.EmailAddress().Equals(confirmEmailAddress.EmailAddress()), ShouldBeTrue)
	So(emailAddressConfirmationFailed.ConfirmationHash().Equals(confirmEmailAddress.ConfirmationHash()), ShouldBeTrue)
	So(emailAddressConfirmationFailed.StreamVersion(), ShouldEqual, 2)
}
