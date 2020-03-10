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

func TestChangeCustomerEmailAddress(t *testing.T) {
	Convey("Prepare test artifacts", t, func() {
		var err error

		customerID := values.GenerateCustomerID()
		emailAddress := values.RebuildEmailAddress("kevin@ball.com")
		confirmationHash := values.GenerateConfirmationHash(emailAddress.EmailAddress())
		personName := values.RebuildPersonName("Kevin", "Ball")
		changedEmailAddress := values.RebuildEmailAddress("latoya@ball.net")

		registered := events.CustomerWasRegistered(
			customerID,
			emailAddress,
			confirmationHash,
			personName,
			1,
		)

		emailAddressConfirmed := events.CustomerEmailAddressWasConfirmed(
			customerID,
			emailAddress,
			2,
		)

		changeEmailAddress, err := commands.NewChangeEmailAddress(
			customerID.ID(),
			changedEmailAddress.EmailAddress(),
		)
		So(err, ShouldBeNil)

		changedConfirmationHash := changeEmailAddress.ConfirmationHash()

		confirmEmailAddress, err := commands.NewConfirmEmailAddress(
			customerID.ID(),
			changedConfirmationHash.Hash(),
		)
		So(err, ShouldBeNil)

		Convey("\nSCENARIO 1: Change a Customer's emailAddress", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{registered}

				Convey("When ChangeCustomerEmailAddress", func() {
					recordedEvents := domain.ChangeCustomerEmailAddress(eventStream, changeEmailAddress)

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
				eventStream := es.DomainEvents{registered}

				Convey("When ChangeCustomerEmailAddress", func() {
					changeEmailAddress, err = commands.NewChangeEmailAddress(
						customerID.ID(),
						emailAddress.EmailAddress(),
					)
					So(err, ShouldBeNil)

					recordedEvents := domain.ChangeCustomerEmailAddress(eventStream, changeEmailAddress)

					Convey("Then no event", func() {
						So(recordedEvents, ShouldBeEmpty)
					})
				})
			})
		})

		Convey("\nSCENARIO 3: Try to change a Customer's emailAddress to the value it was already changed to", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{registered}

				Convey("and CustomerEmailAddressChanged", func() {
					emailAddressChanged := events.CustomerEmailAddressWasChanged(
						customerID,
						changedEmailAddress,
						changedConfirmationHash,
						2,
					)

					eventStream = append(eventStream, emailAddressChanged)

					Convey("When ChangeCustomerEmailAddress", func() {
						recordedEvents := domain.ChangeCustomerEmailAddress(eventStream, changeEmailAddress)

						Convey("Then no event", func() {
							So(recordedEvents, ShouldBeEmpty)
						})
					})
				})
			})
		})

		Convey("\nSCENARIO 4: Confirm a Customer's changed emailAddress, after the original emailAddress was confirmed", func() {
			Convey("Given CustomerRegistered", func() {
				eventStream := es.DomainEvents{registered}

				Convey("and CustomerEmailAddressConfirmed", func() {
					eventStream = append(eventStream, emailAddressConfirmed)

					Convey("and CustomerEmailAddressChanged", func() {
						emailAddressChanged := events.CustomerEmailAddressWasChanged(
							customerID,
							changedEmailAddress,
							changedConfirmationHash,
							3,
						)

						eventStream = append(eventStream, emailAddressChanged)

						Convey("When ConfirmEmailAddress", func() {
							recordedEvents := domain.ConfirmEmailAddress(eventStream, confirmEmailAddress)

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
	})
}
