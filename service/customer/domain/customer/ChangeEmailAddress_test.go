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

func TestChangeEmailAddress(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		var err error
		var recordedEvents es.DomainEvents

		customerID := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("kevin@ball.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("Kevin", "Ball")
		changedEmailAddress := values.RebuildEmailAddress("latoya@ball.net")

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

		changeEmailAddress, err := commands.BuildChangeCustomerEmailAddress(
			customerID.ID(),
			changedEmailAddress.EmailAddress(),
		)
		So(err, ShouldBeNil)

		changedConfirmationHash := changeEmailAddress.ConfirmationHash()

		confirmEmailAddress, err := commands.BuildConfirmCustomerEmailAddress(
			customerID.ID(),
			changedConfirmationHash.Hash(),
		)
		So(err, ShouldBeNil)

		Convey("\nSCENARIO 1: Change a Customer's emailAddress", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("When ChangeCustomerEmailAddress", func() {
					recordedEvents, err = customer.ChangeEmailAddress(eventStream, changeEmailAddress)
					So(err, ShouldBeNil)

					Convey("Then CustomerEmailAddressChanged", func() {
						So(recordedEvents, ShouldHaveLength, 1)
						emailAddressChanged, ok := recordedEvents[0].(events.CustomerEmailAddressChanged)
						So(ok, ShouldBeTrue)
						So(emailAddressChanged, ShouldNotBeNil)
						So(emailAddressChanged.CustomerID().Equals(customerID), ShouldBeTrue)
						So(emailAddressChanged.EmailAddress().Equals(changedEmailAddress), ShouldBeTrue)
						So(emailAddressChanged.ConfirmationHash().Equals(changedConfirmationHash), ShouldBeTrue)
						So(emailAddressChanged.StreamVersion(), ShouldEqual, 2)
					})
				})
			})
		})

		Convey("\nSCENARIO 2: Try to change a Customer's emailAddress to the value he registered with", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("When ChangeCustomerEmailAddress", func() {
					changeEmailAddress, err = commands.BuildChangeCustomerEmailAddress(
						customerID.ID(),
						emailAddress.EmailAddress(),
					)
					So(err, ShouldBeNil)

					recordedEvents, err = customer.ChangeEmailAddress(eventStream, changeEmailAddress)
					So(err, ShouldBeNil)

					Convey("Then no event", func() {
						So(recordedEvents, ShouldBeEmpty)
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Try to change a Customer's emailAddress to the value it was already changed to", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("and CustomerEmailAddressChanged", func() {
					emailAddressChanged := events.CustomerEmailAddressWasChanged(
						customerID,
						changedEmailAddress,
						changedConfirmationHash,
						2,
					)

					eventStream = append(eventStream, emailAddressChanged)

					Convey("When ChangeCustomerEmailAddress", func() {
						recordedEvents, err = customer.ChangeEmailAddress(eventStream, changeEmailAddress)
						So(err, ShouldBeNil)

						Convey("Then no event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 4: Confirm a Customer's changed emailAddress, after the original emailAddress was confirmed", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("and CustomerEmailAddressConfirmed", func() {
					eventStream = append(eventStream, customerEmailAddressWasConfirmed)

					Convey("and CustomerEmailAddressChanged", func() {
						emailAddressChanged := events.CustomerEmailAddressWasChanged(
							customerID,
							changedEmailAddress,
							changedConfirmationHash,
							3,
						)

						eventStream = append(eventStream, emailAddressChanged)

						Convey("When ConfirmCustomerEmailAddress", func() {
							recordedEvents, err = customer.ConfirmEmailAddress(eventStream, confirmEmailAddress)
							So(err, ShouldBeNil)

							Convey("Then CustomerEmailAddressConfirmed", func() {
								So(recordedEvents, ShouldHaveLength, 1)
								emailAddressConfirmed, ok := recordedEvents[0].(events.CustomerEmailAddressConfirmed)
								So(ok, ShouldBeTrue)
								So(emailAddressConfirmed.CustomerID().Equals(customerID), ShouldBeTrue)
								So(emailAddressConfirmed.EmailAddress().Equals(changedEmailAddress), ShouldBeTrue)
								So(emailAddressConfirmed.StreamVersion(), ShouldEqual, 4)
							})
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 5: Try to change a Customer's emailAddress when the account was deleted", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{customerWasRegistered}

				Convey("Given CustomerDeleted", func() {
					eventStream = append(
						eventStream,
						events.CustomerWasDeleted(customerID, 2),
					)

					Convey("When ChangeCustomerEmailAddress", func() {
						_, err := customer.ChangeEmailAddress(eventStream, changeEmailAddress)

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
