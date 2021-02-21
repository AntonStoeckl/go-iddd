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
		emailAddress, err := value.BuildUnconfirmedEmailAddress("kevin@ball.com")
		So(err, ShouldBeNil)
		invalidConfirmationHash := value.RebuildConfirmationHash("invalid_hash")
		personName, err := value.BuildPersonName("Kevin", "Ball")
		So(err, ShouldBeNil)

		command := domain.BuildConfirmCustomerEmailAddress(customerID, emailAddress.ConfirmationHash())
		commandWithInvalidHash := domain.BuildConfirmCustomerEmailAddress(customerID, invalidConfirmationHash)

		customerRegistered := domain.BuildCustomerRegistered(
			customerID,
			emailAddress,
			personName,
			es.GenerateMessageID(),
			1,
		)

		customerEmailAddressConfirmed := domain.BuildCustomerEmailAddressConfirmed(
			customerID,
			value.ToConfirmedEmailAddress(emailAddress),
			es.GenerateMessageID(),
			2,
		)

		customerDeleted := domain.BuildCustomerDeleted(
			customerID,
			es.GenerateMessageID(),
			2,
		)

		Convey("\nSCENARIO 1: ConfirmEmailAddress a Customer's emailAddress with the right confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("When ConfirmCustomerEmailAddress", func() {
					recordedEvents, err = customer.ConfirmEmailAddress(eventStream, command)
					So(err, ShouldBeNil)

					Convey("Then CustomerEmailAddressConfirmed", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						event, ok := recordedEvents[0].(domain.CustomerEmailAddressConfirmed)
						So(ok, ShouldBeTrue)
						So(event.CustomerID().Equals(customerID), ShouldBeTrue)
						So(event.EmailAddress().Equals(emailAddress), ShouldBeTrue)
						So(event.Meta().CausationID(), ShouldEqual, command.MessageID().String())
						So(event.Meta().MessageID(), ShouldNotBeEmpty)
						So(event.Meta().StreamVersion(), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: ConfirmEmailAddress a Customer's emailAddress with a wrong confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("When ConfirmCustomerEmailAddress", func() {
					recordedEvents, err = customer.ConfirmEmailAddress(eventStream, commandWithInvalidHash)
					So(err, ShouldBeNil)

					Convey("Then CustomerEmailAddressConfirmationFailed", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						event, ok := recordedEvents[0].(domain.CustomerEmailAddressConfirmationFailed)
						So(ok, ShouldBeTrue)
						So(event.CustomerID().Equals(customerID), ShouldBeTrue)
						So(event.ConfirmationHash().Equals(invalidConfirmationHash), ShouldBeTrue)
						So(event.IsFailureEvent(), ShouldBeTrue)
						So(event.FailureReason(), ShouldBeError)
						So(event.Meta().CausationID(), ShouldEqual, commandWithInvalidHash.MessageID().String())
						So(event.Meta().MessageID(), ShouldNotBeEmpty)
						So(event.Meta().StreamVersion(), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Try to confirm a Customer's emailAddress again with the right confirmationHash", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("and CustomerEmailAddressConfirmed", func() {
					eventStream = append(eventStream, customerEmailAddressConfirmed)

					Convey("When ConfirmCustomerEmailAddress", func() {
						recordedEvents, err = customer.ConfirmEmailAddress(eventStream, command)
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
				eventStream := es.EventStream{customerRegistered}

				Convey("and CustomerEmailAddressConfirmed", func() {
					eventStream = append(eventStream, customerEmailAddressConfirmed)

					Convey("When ConfirmCustomerEmailAddress", func() {
						recordedEvents, err = customer.ConfirmEmailAddress(eventStream, commandWithInvalidHash)
						So(err, ShouldBeNil)

						Convey("Then no event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 6: Try to confirm a Customer's emailAddress when the account was deleted", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("Given CustomerDeleted", func() {
					eventStream = append(eventStream, customerDeleted)

					Convey("When ConfirmCustomerEmailAddress", func() {
						_, err := customer.ConfirmEmailAddress(eventStream, command)

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

func TestConfirmEmailAddressAfterItWasChanged(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		var err error
		var recordedEvents es.RecordedEvents

		customerID := value.GenerateCustomerID()
		emailAddress, err := value.BuildUnconfirmedEmailAddress("kevin@ball.com")
		So(err, ShouldBeNil)
		changedEmailAddress, err := value.BuildUnconfirmedEmailAddress("latoya@ball.net")
		So(err, ShouldBeNil)
		personName, err := value.BuildPersonName("Kevin", "Ball")
		So(err, ShouldBeNil)

		command := domain.BuildConfirmCustomerEmailAddress(customerID, changedEmailAddress.ConfirmationHash())

		customerRegistered := domain.BuildCustomerRegistered(
			customerID,
			emailAddress,
			personName,
			es.GenerateMessageID(),
			1,
		)

		customerEmailAddressConfirmed := domain.BuildCustomerEmailAddressConfirmed(
			customerID,
			value.ToConfirmedEmailAddress(emailAddress),
			es.GenerateMessageID(),
			2,
		)

		customerEmailAddressChanged := domain.BuildCustomerEmailAddressChanged(
			customerID,
			changedEmailAddress,
			es.GenerateMessageID(),
			3,
		)

		Convey("\nSCENARIO 1: ConfirmEmailAddress a Customer's changed emailAddress, after the original emailAddress was confirmed", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.EventStream{customerRegistered}

				Convey("and CustomerEmailAddressConfirmed", func() {
					eventStream = append(eventStream, customerEmailAddressConfirmed)

					Convey("and CustomerEmailAddressChanged", func() {
						eventStream = append(eventStream, customerEmailAddressChanged)

						Convey("When ConfirmCustomerEmailAddress", func() {
							recordedEvents, err = customer.ConfirmEmailAddress(eventStream, command)
							So(err, ShouldBeNil)

							Convey("Then CustomerEmailAddressConfirmed", func() {
								So(recordedEvents, ShouldHaveLength, 1)
								event, ok := recordedEvents[0].(domain.CustomerEmailAddressConfirmed)
								So(ok, ShouldBeTrue)
								So(event.CustomerID().Equals(customerID), ShouldBeTrue)
								So(event.EmailAddress().Equals(changedEmailAddress), ShouldBeTrue)
								So(event.IsFailureEvent(), ShouldBeFalse)
								So(event.FailureReason(), ShouldBeNil)
								So(event.Meta().CausationID(), ShouldEqual, command.MessageID().String())
								So(event.Meta().MessageID(), ShouldNotBeEmpty)
								So(event.Meta().StreamVersion(), ShouldEqual, 4)
							})
						})
					})
				})
			})
		})
	})
}
